package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
	"time"
)

const (
	defaultExpirationTime = time.Hour
)

// Client used to make requests to redis
type Client struct {
	*redis.Client
	ttl       time.Duration
	namespace string
}

// Param is an optional param for redis client.
type Param func(*Client)

// WithTTL used to set keys expiration time.
func WithTTL(t time.Duration) Param {
	return func(c *Client) {
		c.ttl = t
	}
}

var redisClient *Client
var ctx = context.Background()

// NewClient is a client constructor.
func NewClient(connectionURL, password, namespace string) *Client {
	log.Info("connecting to redis client")

	c := redis.NewClient(&redis.Options{
		Addr:        connectionURL,
		Password:    password, // no password set
		DB:          0,
		DialTimeout: 15 * time.Second,
		MaxRetries:  10, // use default DB
	})

	// Test redis connection
	if _, err := c.Ping(ctx).Result(); err != nil {
		log.Panic("unable to connect to redis: %s", err)
	}

	log.Info("connected to redis client")
	client := &Client{
		Client:    c,
		ttl:       defaultExpirationTime,
		namespace: namespace,
	}

	setRedisClient(client)

	return client
}

func setRedisClient(client *Client) {
	redisClient = client
}

func RedisClient() *Client {
	return redisClient
}

func (c *Client) Ping() error {
	_, err := c.Client.Ping(ctx).Result()
	return err
}
