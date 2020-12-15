package cache

import (
	"encoding/json"
	"time"

	"github.com/cmd-ctrl-q/go-mux-crash-course/entity"
	"github.com/go-redis/redis"
)

/*
this caching layer will be included within the post controller,
so to store the post into memory using redis,
which will help us improve the performance of the api
*/

type redisCache struct {
	host    string
	db      int           // 0-15
	expires time.Duration // expiration in seconds
}

func NewRedisCache(host string, db int, exp time.Duration) PostCache {
	return &redisCache{
		host:    host,
		db:      db,
		expires: exp,
	}
}

// private method - creates a new redis client
func (cache *redisCache) getClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     cache.host,
		Password: "",
		DB:       cache.db,
	})
}

func (cache *redisCache) Set(key string, value *entity.Post) {
	// asscoaite the json representation of the post to the key
	client := cache.getClient()

	json, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}

	// set the key
	client.Set(key, json, cache.expires*time.Second)
}

func (cache *redisCache) Get(key string) *entity.Post {
	client := cache.getClient()

	// get the value of the key
	val, err := client.Get(key).Result()
	if err != nil {
		return nil
	}

	post := entity.Post{}
	// unmarshal the json and store in a new post
	err = json.Unmarshal([]byte(val), &post)
	if err != nil {
		panic(err)
	}

	return &post
}
