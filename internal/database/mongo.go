package database

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
	"mizito/internal/env"
	"mizito/pkg/models"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoHandler interface {
	StoreMessage(payload []byte)

	//GetMessagesSince takes the userid and a date, returning all the messages sent after the date
	GetMessagesSince(context context.Context, to uint, sinceDate time.Time) error
}

type mongoHandler struct {
	client      *mongo.Client
	redis       RedisHandler
	cfg         *env.Config
	messageChan chan []byte
	messagesLen int
}

func NewMongoHandler(env *env.Config, Redis RedisHandler) MongoHandler {
	var mongoDB mongoHandler
	if client, err := mongo.Connect(options.Client().ApplyURI(env.MongoDBHost)); err != nil {
		panic(fmt.Sprintf("failed to establish connection to mongodb, err: %s", err))
	} else {
		mongoDB = mongoHandler{
			client:      client,
			redis:       Redis,
			messageChan: make(chan []byte, 100),
			cfg:         env,
		}
	}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	if err := mongoDB.client.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}
	Migrate(ctx, mongoDB, env.MongoDatabase, env.MongoCollection)

	go mongoDB.ProcessMessages()

	return &mongoDB
}

func Migrate(ctx context.Context, mongoDB mongoHandler, dbname string, collectionName string) {
	db := mongoDB.client.Database(dbname)
	names, err := db.ListCollectionNames(context.Background(), bson.D{})
	if err != nil {
		panic(err)
	}

	var found bool = false

	for _, name := range names {
		if name == collectionName {
			found = true
		}
	}

	if !found {
		if err := db.CreateCollection(ctx, collectionName); err != nil {
			panic(err)
		}
	}
}

func (mh *mongoHandler) StoreMessage(event []byte) {
	mh.messageChan <- event
}

func (mh *mongoHandler) ProcessMessages() {
	db := mh.client.Database(mh.cfg.MongoDatabase)
	coll := db.Collection(mh.cfg.MongoCollection)

	for message := range mh.messageChan {
		ctx, _ := context.WithTimeout(context.Background(), time.Second*5)

		var doc bson.D

		fmt.Println(string(message))

		if err := json.Unmarshal(message, &doc); err != nil {
			// error and log
			continue
		}

		fmt.Println(doc)

		res, err := coll.InsertOne(ctx, doc)
		if err != nil {
			//handle errors here
			fmt.Println("error occurred while inserting item into db", err)
			continue
		}

		if !res.Acknowledged {
			fmt.Println("ack not received")
			// handle no acknowledge received error
		}

		//if insertion is acknowledged, send a corresponding event to redis
		mh.redis.AddToPublishChan(message)
	}

}

func (mh *mongoHandler) GetMessagesSince(ctx context.Context, to uint, sinceDate time.Time) error {
	db := mh.client.Database(mh.cfg.MongoDatabase)
	coll := db.Collection(mh.cfg.MongoCollection)

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
