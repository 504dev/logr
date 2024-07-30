package types

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"time"
)

type SockSessionStore interface {
	Set(key string, value *SockSession) error
	Get(key string) (*SockSession, error)
	Del(key string) error
}

type MemorySessionStore struct{}

func (s MemorySessionStore) Set(key string, value *SockSession) error {
	return nil
}
func (s MemorySessionStore) Get(key string) (*SockSession, error) {
	return nil, nil
}
func (s MemorySessionStore) Del(key string) error {
	return nil
}

func NewRedisSessionStore(addr string, db int, ttl time.Duration) (*RedisSessionStore, error) {
	opt, err := redis.ParseURL(addr)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opt)
	ctx := context.Background()

	_, err = client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return &RedisSessionStore{
		client: client,
		ctx:    ctx,
		ttl:    ttl,
	}, nil
}

type RedisSessionStore struct {
	client *redis.Client
	ctx    context.Context
	ttl    time.Duration
}

func (store *RedisSessionStore) Set(key string, value *SockSession) error {
	json, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return store.client.Set(store.ctx, key, json, store.ttl).Err()
}

func (store *RedisSessionStore) Get(key string) (*SockSession, error) {
	val, err := store.client.Get(store.ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var session SockSession
	err = json.Unmarshal([]byte(val), &session)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (store *RedisSessionStore) Del(key string) error {
	return store.client.Del(store.ctx, key).Err()
}
