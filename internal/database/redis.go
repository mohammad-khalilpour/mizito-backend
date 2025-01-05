package database

import (
	"context"
	"github.com/redis/go-redis/v9"
)

type RedisHandler interface {
	AddToPublishChan(message []byte)
	GetSubscribeChan() chan []byte
}

type redisHandler struct {
	client        redis.Client
	publishChan   chan *redis.Message
	SubscribeChan chan []byte
}

func NewRedisHandler() RedisHandler {
	//TODO

	return nil
}

func (rm *redisHandler) Subscribe(ctx context.Context) {
	subscriber := rm.client.Subscribe(ctx, "messages")

	for {
		msg, err := subscriber.ReceiveMessage(ctx)
		if err != nil {
			// log error
		}
		rm.SubscribeChan <- []byte(msg.Payload)
	}
}

func (rm *redisHandler) Publish(ctx context.Context) {
	for event := range rm.publishChan {
		if inf := rm.client.Publish(ctx, event.Channel, event.Payload); inf.Err() != nil {
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
