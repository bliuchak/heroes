package db

import (
	"strconv"
	"strings"

	"github.com/bliuchak/heroes/internal/storage"
	"github.com/go-redis/redis"
)

const (
	heroPrefix = "hero"
)

type Redis struct {
	Client redis.Client
}

func NewRedis(host, password string, port int) (*Redis, error) {
	opt := redis.Options{
		Addr:     host + ":" + strconv.Itoa(port),
		Password: password,
		DB:       0,
	}
	c := redis.NewClient(&opt)

	_, err := c.Ping().Result()
	if err != nil {
		return &Redis{}, err
	}

	return &Redis{Client: *c}, nil
}

func (r *Redis) Status() (string, error) {
	return r.Client.Ping().Result()
}

func (r *Redis) GetHeroes() ([]storage.Hero, error) {
	var heroes []storage.Hero

	iter := r.Client.Scan(0, heroPrefix+".*", 100).Iterator()
	for iter.Next() {
		id := strings.Split(iter.Val(), ".")
		name, err := r.Client.Get(iter.Val()).Result()
		if err != nil {
			return []storage.Hero{}, err
		}
		heroes = append(heroes, storage.Hero{ID: id[1], Name: name})
	}
	if err := iter.Err(); err != nil {
		return []storage.Hero{}, err
	}
	return heroes, nil
}

func (r *Redis) GetHero(id string) (storage.Hero, error) {
	res, err := r.Client.Get(heroPrefix + "." + id).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			return storage.Hero{}, storage.NewErrHeroNotExist("hero not exist")
		}
		return storage.Hero{}, err
	}

	return storage.Hero{ID: id, Name: res}, nil
}

func (r *Redis) CreateHero(id, name string) error {
	_, err := r.Client.Set(heroPrefix+"."+id, name, 0).Result()
	if err != nil {
		return err
	}

	return nil
}

func (r *Redis) DeleteHero(id string) error {
	res, err := r.Client.Del(heroPrefix + "." + id).Result()
	if err != nil {
		return err
	}

	if res == 0 {
		return storage.NewErrNothingToDelete("nothing to delete")
	}

	return nil
}
