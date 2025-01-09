package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/v2/bson"
	"mizito/internal/database"
	"mizito/internal/env"
	"mizito/pkg/models"
	"mizito/pkg/models/dtos"
	messagedto "mizito/pkg/models/dtos/message"
	"time"
)

type MessageStoreRepository interface {
	StoreMessage(message *messagedto.Message) error
	GetMessagesSince(ctx context.Context, to uint, sinceDate time.Time) error
}

type MessageChannelRepository interface {
	PublishMsg(event []byte)
	PublishEvent(event dtos.Event)
	SubscribeEvent() <-chan dtos.Event
}

type MessageRepository interface {
	MessageChannelRepository
	MessageStoreRepository
}

type messageRepository struct {
	redis        database.RedisHandler
	mongo        database.MongoHandler
	mongoChan    chan []byte
	redisPubChan chan dtos.Event
	redisSubChan chan dtos.Event
	messageLen   int
	cfg          *env.Config
}

func NewMessageRepository(redis *database.RedisHandler, mongo *database.MongoHandler, env *env.Config) MessageRepository {
	msgRepo := messageRepository{
		redis:        *redis,
		mongo:        *mongo,
		mongoChan:    make(chan []byte, 100),
		redisPubChan: make(chan dtos.Event, 100),
		redisSubChan: make(chan dtos.Event, 100),
		cfg:          env,
		messageLen:   100,
	}

	go msgRepo.ProcessMessage()

	go msgRepo.StoreEventMessage()

	go msgRepo.SubscribeMessages()

	return &msgRepo

}

func (mr *messageRepository) GetMessagesSince(ctx context.Context, to uint, sinceDate time.Time) error {
	db := mr.mongo.Client.Database(mr.cfg.MongoDatabase)
	coll := db.Collection(mr.cfg.MongoCollection)

	c, err := coll.Find(ctx, bson.D{{"created_at", bson.D{{"$gt", sinceDate}}},
		{"to", bson.D{{"$contains", to}}}})

	if err != nil {
		return fmt.Errorf("failed to ")
	}

	var messages []models.Message

	if err := c.All(ctx, &messages); err != nil {
		return fmt.Errorf("failed to cast documents as message type, err : %s", err.Error())
	}

	return nil
}

func (mr *messageRepository) PublishMsg(event []byte) {
	mr.mongoChan <- event
}

func (mr *messageRepository) PublishEvent(event dtos.Event) {
	mr.redisPubChan <- event
}

func (rm *messageRepository) SubscribeEvent() <-chan dtos.Event {
	return rm.redisSubChan
}

func (mr *messageRepository) ProcessMessage() {

	for message := range mr.mongoChan {

		var event dtos.Event
		if err := json.Unmarshal(message, &event); err != nil {
			fmt.Printf("failed to parse the event into Event, err :%s\n", err.Error())
			// log for error
		}
		if err := mr.StoreMessage(&event.Payload); err != nil {
			fmt.Printf("failed to insert message into db, err: %s\n", err.Error())
		}
		mr.PublishEvent(event)
	}

}

func (rm *messageRepository) SubscribeMessages() {
	subscriber := rm.redis.Client.Subscribe(context.Background(), "messages")

	for {
		ctx, _ := context.WithTimeout(context.Background(), time.Second*2)
		event, err := subscriber.ReceiveMessage(ctx)
		if err != nil {
			// log error
			continue
		}

		var msg dtos.Event

		if err := json.Unmarshal([]byte(event.Payload), &msg); err != nil {
			// log for encountering error
			continue
		}

		rm.redisSubChan <- msg
	}
}

func (mr *messageRepository) StoreEventMessage() {
	for event := range mr.redisPubChan {
		redisPayload, err := json.Marshal(event)

		if err != nil {
			//log the malicious event
		}

		if err := mr.redis.Publish(redisPayload, "messages"); err != nil {
			fmt.Println("failed to store the event inside the database")
		}
	}
}

func (mr *messageRepository) StoreMessage(message *messagedto.Message) error {
	db := mr.mongo.Client.Database(mr.cfg.MongoDatabase)
	coll := db.Collection(mr.cfg.MongoCollection)

	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)

	res, err := coll.InsertOne(ctx, message)
	if err != nil {
		//handle errors here
		return fmt.Errorf("error occurred while inserting item into db: %s", err.Error())
	}

	if !res.Acknowledged {
		fmt.Println("ack not received")
		// handle no acknowledge received error
	}

	return nil

}
