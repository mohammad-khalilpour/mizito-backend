package database

import (
	"context"

	"github.com/redis/go-redis/v9"
)



type RedisManager struct {
	client	redis.Client
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