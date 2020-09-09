package accounts

import "github.com/go-redis/redis"

type Publisher struct {
	rdb *redis.Client
}

func NewPublisher(rdb *redis.Client) *Publisher {
	return &Publisher{rdb}
}

func (publisher *Publisher) Publish(channel string, message string) error {
	return publisher.rdb.Publish(channel, message).Err()
}
