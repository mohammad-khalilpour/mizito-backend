package database

import (
	"context"
	"encoding/json"
	"fmt"
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
	GetMessagesSince(to uint, sinceDate time.Time)
}

type mongoHandler struct {
	client      *mongo.Client
	cfg         *env.Config
	messageChan chan bson.M
	messagesLen int
}

func (mh mongoHandler) StoreMessage(event []byte) {
	var m map[string]any
	if err := json.Unmarshal(event, &m); err != nil {
		return
	}
	mh.messageChan <- m
}

func (mh *mongoHandler) ProcessMessages() {

	var messages []bson.M

	db := mh.client.Database(mh.cfg.MongoDatabase)
	coll := db.Collection(mh.cfg.MongoCollection)

	for message := range mh.messageChan {

		messages = append(messages, message)

		if len(messages) == mh.messagesLen {
			go func(messages []bson.M) {
				ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
				res, errs := coll.InsertMany(ctx, messages)
				if errs != nil {
					//handle errors here
				}

				if !res.Acknowledged {
					// handle no acknowledge received error
				}
			}(messages)
		}
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

var mongoDB mongoHandler

func Connect() {
	if mongoDB.client != nil {
		return
	}
	if client, err := mongo.Connect(options.Client().ApplyURI("mongodb://localhost:27017")); err != nil {
		panic(fmt.Sprintf("failed to establish connection to mongodb, err: %s", err))
	} else {
		mongoDB = mongoHandler{client: client}
	}
}

func Migrate(ctx context.Context, dbname string, collectionName string) {
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
