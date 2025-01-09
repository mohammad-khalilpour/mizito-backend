package database

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
	"mizito/internal/env"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoHandler struct {
	Client *mongo.Client
	Redis  RedisHandler
	Cfg    *env.Config
}

func NewMongoHandler(env *env.Config) *MongoHandler {
	var mongoDB MongoHandler
	if client, err := mongo.Connect(options.Client().ApplyURI(env.MongoDBHost)); err != nil {
		panic(fmt.Sprintf("failed to establish connection to mongodb, err: %s", err))
	} else {
		mongoDB = MongoHandler{
			Client: client,
			Cfg:    env,
		}
	}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	if err := mongoDB.Client.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}
	Migrate(ctx, mongoDB, env.MongoDatabase, env.MongoCollection)

	return &mongoDB
}

func Migrate(ctx context.Context, mongoDB MongoHandler, dbname string, collectionName string) {
	db := mongoDB.Client.Database(dbname)
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
