/*
Package redis adds built-in Redis support with our common configuration.
This will ensure that Go services using "github.com/skuid/spec" are in compliance with our infrastructure expectations.

The package provides two utilities for using Redis.
NewStandardRedisClient generates a full-featured Redis client, with preconfigured connection settings from our spec.
NewStandardRedisCache generates a more limited Redis cache interface whcih will automatically marshal
and unmarshal into msgpack, also preconfigured with connection settings from our spec.
*/
package redis

import (
	"os"

	"github.com/go-redis/cache"
	"github.com/go-redis/redis"
	"gopkg.in/vmihailenco/msgpack.v2"
)

var redisHost string

func init() {
	redisHost = os.Getenv("REDIS_HOST")
}

// NewStandardRedisClient generates a preconfigured redis client according to our spec.
// Accepts an *redis.Options object, and overrides the Addr field to use `$REDIS_HOST:6379` instead
func NewStandardRedisClient(options *redis.Options) *redis.Client {
	options.Addr = redisHost + ":6379"
	return redis.NewClient(options)
}

// NewStandardRedisCache generates a preconfigured redis cache according to our spec, using msgpack for serialization format.
// Accepts an *redis.Options object, and overrides the Addr field to use `$REDIS_HOST:6379` instead
func NewStandardRedisCache(options *redis.Options) *cache.Codec {
	return &cache.Codec{
		Redis: NewStandardRedisClient(options),
		Marshal: func(v interface{}) ([]byte, error) {
			return msgpack.Marshal(v)
		},
		Unmarshal: func(b []byte, v interface{}) error {
			return msgpack.Unmarshal(b, v)
		},
	}
}
