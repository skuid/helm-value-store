package redis

import (
	"os"
	"testing"

	"github.com/go-redis/cache"
	"github.com/go-redis/redis"
)

func TestClientConnectsAsExpected(t *testing.T) {
	formerHost := os.Getenv("REDIS_HOST")
	os.Setenv("REDIS_HOST", "localhost")

	cases := []struct {
		testKey   string
		testValue string
	}{
		{
			"some test key",
			"some value",
		},
		{
			"skuid",
			"101",
		},
	}

	for _, c := range cases {
		client := NewStandardRedisClient(&redis.Options{
			Password: "",
			DB:       0,
		})
		err := client.Set(c.testKey, c.testValue, 0).Err()
		if err != nil {
			t.Errorf("Failed TestClientConnectsAsExpected: Error in value storage.\n %v", err.Error())
		}
		got, err := client.Get(c.testKey).Result()
		if err != nil {
			t.Errorf("Failed TestClientConnectsAsExpected: Error in value retrieval.\n %v", err.Error())
		}
		if got != c.testValue {
			t.Errorf("Failed TestClientConnectsAsExpected: Key/Value mismatch.\n Expected %#v, got %#v", c.testValue, got)
		}
	}

	os.Setenv("REDIS_HOST", formerHost)
}

func TestCacheConnectsAsExpected(t *testing.T) {
	formerHost := os.Getenv("REDIS_HOST")
	os.Setenv("REDIS_HOST", "localhost")

	type TestObject struct {
		TestString    string
		TestNumerical int
		SubObject     interface{}
	}

	cases := []struct {
		testKey   string
		testValue *TestObject
	}{
		{
			"test key 1",
			&TestObject{
				TestString:    "test string value",
				TestNumerical: 1990,
				SubObject:     "String SubObject",
			},
		},
	}

	for _, c := range cases {
		testCache := NewStandardRedisCache(&redis.Options{
			Password: "",
			DB:       0,
		})
		err := testCache.Set(&cache.Item{
			Key:    c.testKey,
			Object: c.testValue,
		})
		if err != nil {
			t.Errorf("Failed TestClientConnectsAsExpected: Error in value storage.\n %v", err.Error())
		}
		var got TestObject
		err = testCache.Get(c.testKey, &got)
		if err != nil {
			t.Errorf("Failed TestClientConnectsAsExpected: Error in value retrieval.\n %v", err.Error())
		}
		if got != *c.testValue {
			t.Errorf("Failed TestClientConnectsAsExpected: Key/Value mismatch.\n Expected %#v, got %#v", c.testValue, got)
		}
	}

	os.Setenv("REDIS_HOST", formerHost)
}
