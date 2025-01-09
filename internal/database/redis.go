package database

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"mizito/internal/env"
	"mizito/pkg/models/dtos"
)

type RedisHandler struct {
	Client        *redis.Client
	publishChan   chan *redis.Message
	SubscribeChan chan dtos.Event
}

func NewRedisHandler(env *env.Config) *RedisHandler {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", env.RedisHost, env.RedisPort),
		Username: env.RedisUsername,
		Password: env.RedisPassword,
	})

	//ctx, _ := context.WithTimeout(context.Background(), time.Second*4)
	if err := client.Ping(context.Background()).Err(); err != nil {
		panic(fmt.Sprintf("failed to connect to redis db, err : %s", err.Error()))
	}

	r := &RedisHandler{
		Client:        client,
		publishChan:   make(chan *redis.Message, 100),
		SubscribeChan: make(chan dtos.Event, 100),
	}

	return r
}

func (rm *RedisHandler) Publish(event []byte, channel string) error {
	if inf := rm.Client.Publish(context.Background(), channel, event); inf.Err() != nil {
		return fmt.Errorf("failed to publish message into channel, err : %s", inf.Err().Error())
	}
	return nil
}
