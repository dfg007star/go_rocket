package session

import (
	"context"
	"errors"
	"fmt"
	"time"

	redigo "github.com/gomodule/redigo/redis"
)

type repository struct {
	cache cache.RedisClient
}

func NewRepository(cache cache.RedisClient) *repository {
	return &repository{
		cache: cache,
	}
}
