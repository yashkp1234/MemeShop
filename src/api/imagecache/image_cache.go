package imagecache

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/go-redis/redis"
)

const (
	//ImageCacheTime is how long something lasts in cache
	ImageCacheTime = 30 * time.Minute
	keyName        = "Pictures"
	threshold      = 20
)

//ImageCache object to handle caching of images
type ImageCache struct {
	RedisInstance *redis.Client
}

//ImageCache represents an image cache
var imageCache *ImageCache

//NewImageCacheClient setups creating a image cache connection
func NewImageCacheClient() {
	redisInstance := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	if err := redisInstance.Ping().Err(); err != nil {
		log.Fatal(err)
	}

	imageCache = &ImageCache{redisInstance}
}

//Cancel disconnects with the imagecache
func Cancel() {
	if err := imageCache.RedisInstance.Close(); err != nil {
		log.Fatal(err)
	}
}

//Connect returns the connected redis client
func Connect() *ImageCache {
	return imageCache
}

//AddImage adds an image  to redis
func (c *ImageCache) AddImage(hash interface{}, title interface{}) (string, error) {
	id, err := c.RedisInstance.Do("imgscout.add", "Pictures", hash, title).Result()
	return fmt.Sprint(id), err
}

//SyncImages syncs all recently added images to redis
func (c *ImageCache) SyncImages() error {
	_, err := c.RedisInstance.Do("imgscout.sync", keyName).Result()
	return err
}

//DeleteImage deletes an image from redis
func (c *ImageCache) DeleteImage(id interface{}) error {
	_, err := c.RedisInstance.Do("imgscout.del", keyName, id).Result()
	return err
}

//DeleteAll deltes all the images
func (c *ImageCache) DeleteAll(logging bool) error {
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
func (c *ImageCache) QueryImages(hashVal interface{}) (bool, error) {
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
