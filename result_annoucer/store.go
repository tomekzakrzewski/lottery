package main

import (
	"encoding/json"
	"time"

	"github.com/go-redis/redis"
	"github.com/tomekzakrzewski/lottery/types"
)

type RedisStore struct {
	c *redis.Client
}

func NewRedisStore(c *redis.Client) *RedisStore {
	return &RedisStore{
		c: c,
	}
}

func (s *RedisStore) Insert(result *types.ResultResponse) error {
	resultJson, _ := json.Marshal(result)
	err := s.c.Set(result.Hash, resultJson, 24*time.Hour).Err()
	if err != nil {
		return err
	}

	return nil
}

func (s *RedisStore) Find(hash string) (*types.ResultResponse, error) {
	var result types.ResultResponse
	resultJson, err := s.c.Get(hash).Result()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(resultJson), &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
