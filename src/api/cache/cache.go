package cache

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/go-redis/redis"
)

const (
	//CacheTime is how long something lasts in cache
	CacheTime = 30 * time.Minute
	keyName   = "Pictures"
	threshold = 20
)

//Cache object to handle caching of images
type Cache struct {
	RedisInstance *redis.Client
}

//Cache represents an image cache
var cache *Cache

//NewCacheClient setups creating a image cache connection
func NewCacheClient() {
	redisInstance := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	if err := redisInstance.Ping().Err(); err != nil {
		log.Fatal(err)
	}

	cache = &Cache{redisInstance}
}

//Cancel disconnects with the imagecache
func Cancel() {
	if err := cache.RedisInstance.Close(); err != nil {
		log.Fatal(err)
	}
}

//Connect returns the connected redis client
func Connect() *Cache {
	return cache
}

//AddImage adds an image  to redis
func (c *Cache) AddImage(hash interface{}, title interface{}) (string, error) {
	id, err := c.RedisInstance.Do("imgscout.add", "Pictures", hash, title).Result()
	return fmt.Sprint(id), err
}

//SyncImages syncs all recently added images to redis
func (c *Cache) SyncImages() error {
	_, err := c.RedisInstance.Do("imgscout.sync", keyName).Result()
	return err
}

//DeleteImage deletes an image from redis
func (c *Cache) DeleteImage(id interface{}) error {
	_, err := c.RedisInstance.Do("imgscout.del", keyName, id).Result()
	return err
}

//DeleteAll deltes all the images
func (c *Cache) DeleteAll(logging bool) error {
	var MAXVAL uint64 = 9999999999999999999
	res, err := c.RedisInstance.Do("imgscout.query", keyName, 0, MAXVAL).Result()
	if err != nil {
		return err
	}

	s := fmt.Sprint(res)
	s = strings.ReplaceAll(s, "[", "")
	s = strings.ReplaceAll(s, "]", "")
	t := strings.Split(s, " ")
	for i := 0; i < len(t); i += 3 {
		if logging {
			log.Println(t[i : i+3])
		}
		if err := c.DeleteImage(t[i+1]); err != nil {
			return err
		}
	}
	return nil
}

//QueryImages finds any similar images
func (c *Cache) QueryImages(hashVal interface{}) (bool, error) {
	results, err := c.RedisInstance.Do("imgscout.query", keyName, hashVal, threshold).Result()
	if err != nil {
		if err.Error() == "ERR - no such key" {
			log.Println("First key")
			return true, nil
		}
		return false, err
	}
	return len(fmt.Sprint(results)) == 2, nil
}
