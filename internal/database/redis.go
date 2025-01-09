package database

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"mizito/internal/env"
	"time"
)

type RedisHandler interface {
	AddToPublishChan(message []byte)
	GetSubscribeChan() <-chan []byte
}

type redisHandler struct {
	client        *redis.Client
	publishChan   chan *redis.Message
	SubscribeChan chan []byte
}

func NewRedisHandler(env *env.Config) RedisHandler {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", env.RedisHost, env.RedisPort),
		Username: env.RedisUsername,
		Password: env.RedisPassword,
	})

	if client == nil {
		panic("failed to connect to redis db")
	}

	r := &redisHandler{
		client:        client,
		publishChan:   make(chan *redis.Message, 100),
		SubscribeChan: make(chan []byte, 100),
	}

	go r.Publish()
	go r.Subscribe()

	return r
}

func (rm *redisHandler) Subscribe() {
	subscriber := rm.client.Subscribe(context.Background(), "messages")

	for {
		ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
		msg, err := subscriber.ReceiveMessage(ctx)
		if err != nil {
			// log error
			continue
		}
		rm.SubscribeChan <- []byte(msg.Payload)
	}
}

func (rm *redisHandler) Publish() {
	for event := range rm.publishChan {
		if inf := rm.client.Publish(context.Background(), event.Channel, event.Payload); inf.Err() != nil {
			// log the error
		}
	}
}

func (rm *redisHandler) AddToPublishChan(message []byte) {
	redisMsg := redis.Message{
		Payload: string(message),
		Channel: "messages",
	}
	rm.publishChan <- &redisMsg
}

func (rm *redisHandler) GetSubscribeChan() <-chan []byte {
	return rm.SubscribeChan
}
