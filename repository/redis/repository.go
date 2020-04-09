package redis

import (
	"github.com/alonyb/url_shortener/shortener"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	"log"
)

type redisRepository struct {
	client *redis.Client
}

//buscar builder
func (r redisRepository) Find(code string) (*shortener.Redirect, error) {
	redirect := &shortener.Redirect{}
	key := r.generateKey(code)
	data, err := r.client.HGetAll(key).Result()
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errors.Wrap(shortener.ErrRedirectNotFound, "redis.repository.Find")
	}
	if err != nil {
		log.Println(err)
		return nil, err
	}
	redirect.Code = data["code"]
	redirect.URL = data["url"]
	redirect.CreateAt = data["reated_at"]
	return redirect, nil

}

func (r redisRepository) Store(redirect *shortener.Redirect) error {
	key := r.generateKey(redirect.Code)
	data := map[string]interface{}{
		"code":      redirect.Code,
		"url":       redirect.URL,
		"create_at": redirect.CreateAt,
	}
	_, err := r.client.HMSet(key,data).Result()
	if err != nil {
		return errors.Wrap(err, "repository.Redirect.Store")
	}
	return nil
}

func (r *redisRepository) generateKey(code string) string {
	return fmt.Sprintf("redirect:%s", code)
}

func newRedisClient(redisUrl string) (*redis.Client, error) {
	opts, err := redis.ParseURL(redisUrl)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(opts)
	_, err = client.Ping().Result()
	if err != nil {
		panic(err)
	}
	return client, err
}

func NewRedisRepository(redisUrl string) (shortener.RedirectRepository, error) {
	repo := &redisRepository{}
	client, err := newRedisClient(redisUrl)
	if err != nil {
		return nil, errors.Wrap(err, "repository.redis.newRedisClient")
	}
	repo.client = client
	return repo, nil
}
