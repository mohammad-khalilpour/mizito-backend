package database

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
)


type RedisManager struct {
	client	redis.Client
}


type RedisEvent struct {
	Subject string
	EventType string
	Payload []byte
}


func (rm *RedisManager) Subscribe(ctx context.Context, ch chan<- *redis.Message) {
	subscriber := rm.client.Subscribe(ctx, "messages")


	for {
		msg, err := subscriber.ReceiveMessage(ctx)
		if err != nil {
			// log error
		}
		ch <- msg
	}
}


func (rm *RedisManager) Publish(ctx context.Context, msg []byte) error{


	if inf := rm.client.Publish(ctx, "messages", msg); inf.Err() != nil {
		return inf.Err()
	}
	return nil
}

func (rm *RedisManager) PublishEvent(ctx context.Context, event *RedisEvent) error{
	d, err := json.Marshal(event)
	
	if err != nil {
		return err
	}

	return rm.Publish(ctx, d)
}