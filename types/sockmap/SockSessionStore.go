package sockmap

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-redis/redis/v8"
	"time"
)

type SessionStore interface {
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

type RedisSessionStore struct {
	client  *redis.Client
	ttl     time.Duration
	timeout time.Duration
}

func NewRedisSessionStore(addr string, ttl time.Duration, timeout time.Duration) (*RedisSessionStore, error) {
	opt, err := redis.ParseURL(addr)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opt)

	_, err = client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return &RedisSessionStore{
		client:  client,
		ttl:     ttl,
		timeout: timeout,
	}, nil
}

func (store *RedisSessionStore) Set(key string, value *SockSession) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), store.timeout)
	defer cancel()

	return store.client.Set(ctx, key, bytes, store.ttl).Err()
}

func (store *RedisSessionStore) Get(key string) (*SockSession, error) {
	ctx, cancel := context.WithTimeout(context.Background(), store.timeout)
	defer cancel()

	val, err := store.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
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
	ctx, cancel := context.WithTimeout(context.Background(), store.timeout)
	defer cancel()

	return store.client.Del(ctx, key).Err()
}
