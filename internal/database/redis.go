package database

import (
	"context"
	"github.com/redis/go-redis/v9"
)

type RedisManager struct {
	client        redis.Client
	publishChan   <-chan *redis.Message
	SubscribeChan chan<- *redis.Message
}

func (rm *RedisManager) Subscribe(ctx context.Context) {
	subscriber := rm.client.Subscribe(ctx, "messages")

	for {
		msg, err := subscriber.ReceiveMessage(ctx)
		if err != nil {
			// log error
		}
		rm.SubscribeChan <- msg
	}
}

func (rm *RedisManager) Publish(ctx context.Context) {
	for event := range rm.publishChan {
		if inf := rm.client.Publish(ctx, event.Channel, event.Payload); inf.Err() != nil {
			// log the error
		}
	}
}
