package database

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"mizito/internal/env"
	"mizito/pkg/models/dtos"
)

// RedisHandler manages Redis operations.
type RedisHandler struct {
	Client        *redis.Client
	publishChan   chan *redis.Message
	SubscribeChan chan dtos.Event
}

// NewRedisHandler initializes a new RedisHandler.
func NewRedisHandler(env *env.Config) *RedisHandler {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", env.RedisHost, env.RedisPort),
		Username: env.RedisUsername,
		Password: env.RedisPassword,
		DB:       0, // use default DB
	})

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

// Publish sends a message to a Redis channel.
func (rm *RedisHandler) Publish(event []byte, channel string) error {
	if inf := rm.Client.Publish(context.Background(), channel, event); inf.Err() != nil {
		return fmt.Errorf("failed to publish message into channel, err : %s", inf.Err().Error())
	}
	return nil
}

// SetBlacklistedToken adds a token to the blacklist with a TTL.
func (rm *RedisHandler) SetBlacklistedToken(token string, ttl time.Duration) error {
	return rm.Client.Set(context.Background(), "blacklist:"+token, "true", ttl).Err()
}

// IsTokenBlacklisted checks if a token is blacklisted.
func (rm *RedisHandler) IsTokenBlacklisted(token string) (bool, error) {
	val, err := rm.Client.Get(context.Background(), "blacklist:"+token).Result()
	if err == redis.Nil {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return val == "true", nil
}
