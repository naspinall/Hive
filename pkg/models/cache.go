package models

import (
	"encoding/json"
	"log"
	"time"

	"github.com/go-redis/redis"
)

type CacheService interface {
	Get(key string) (interface{}, error)
	Set(key string, data interface{})
	Remove(key string) error
}

type cacheRedis struct {
	cache *redis.Client
}

func NewRedisCache(addr, password string, db int) (*cacheRedis, error) {
	// Creating redis client
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	// Checking connection to redis.
	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}
	log.Println("Successfully connected to Redis")

	return &cacheRedis{
		cache: client,
	}, nil
}

func (cr *cacheRedis) Get(key string) (interface{}, error) {
	var decodedData interface{}
	data, err := cr.cache.Get(key).Bytes()
	// Do some more error handling here.
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &decodedData)
	if err != nil {
		return nil, err
	}

	return decodedData, nil

}
func (cr *cacheRedis) Set(key string, data interface{}) {

	go func() {
		// Encoding data as JSON.
		encodedData, err := json.Marshal(data)
		if err != nil {
			log.Printf("Bad JSON Marshal for %s", key)
		}

		// Setting JSON string as cache value.
		cr.cache.Set(key, encodedData, time.Hour)

		// Do some more error handling here.
		if err != nil {
			log.Printf("Bad JSON Marshal for %s", key)
		}
	}()
}

func (cr *cacheRedis) Remove(key string) error {
	return cr.cache.Del(key).Err()
}
