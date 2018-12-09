package db

import (
	"strings"

	"github.com/bliuchak/heroes/internal/storage"
	"github.com/mediocregopher/radix/v3"
)

const (
	heroPrefix = "hero"
)

// Redis contains client which operates with storage
type Redis struct {
	client radix.Client
}

// NewRedis returns pointer to Redis structure with filled data
func NewRedis(host, password, port string) (*Redis, error) {
	pool, err := radix.NewPool("tcp", host+":"+port, 10)
	if err != nil {
		return nil, err
	}

	return &Redis{client: pool}, nil
}

// Status checks storage connection status
func (r *Redis) Status() (string, error) {
	var status string
	if err := r.client.Do(radix.Cmd(&status, "PING")); err != nil {
		return status, err
	}

	return status, nil
}

// GetHeroes gets all heroes
func (r *Redis) GetHeroes() ([]storage.Hero, error) {
	var heroes []storage.Hero

	opts := radix.ScanOpts{
		Command: "SCAN",
		Pattern: heroPrefix + ".*",
		Count:   100,
	}
	scanner := radix.NewScanner(r.client, opts)

	var key string
	for scanner.Next(&key) {
		id := strings.Split(key, ".")
		var name string
		if err := r.client.Do(radix.Cmd(&name, "GET", key)); err != nil {
			return []storage.Hero{}, err
		}
		heroes = append(heroes, storage.Hero{ID: id[1], Name: name})
	}

	if err := scanner.Close(); err != nil {
		return []storage.Hero{}, err
	}

	return heroes, nil
}

// GetHero gets hero by ID
func (r *Redis) GetHero(id string) (storage.Hero, error) {
	var exists int
	if err := r.client.Do(radix.Cmd(&exists, "EXISTS", heroPrefix+"."+id)); err != nil {
		return storage.Hero{}, err
	}

	if exists == 0 {
		return storage.Hero{}, storage.NewErrHeroNotExist("hero not exist")
	}

	var name string
	if err := r.client.Do(radix.Cmd(&name, "GET", heroPrefix+"."+id)); err != nil {
		return storage.Hero{}, err
	}

	return storage.Hero{ID: id, Name: name}, nil
}

// CreateHero creates new hero by ID and Name
func (r *Redis) CreateHero(id, name string) error {
	if err := r.client.Do(radix.Cmd(&name, "SET", heroPrefix+"."+id, name)); err != nil {
		return err
	}

	return nil
}

// DeleteHero deletes hero by ID
func (r *Redis) DeleteHero(id string) error {
	if err := r.client.Do(radix.Cmd(nil, "DEL", heroPrefix+"."+id)); err != nil {
		return err
	}

	return nil
}
