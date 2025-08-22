package session

import (
	"fmt"

	"github.com/dfg007star/go_rocket/platform/pkg/cache"
)

const (
	cacheKeyUserSessionPrefix = "user:session:"
	cacheKeyUserSetPrefix     = "user:set:"
)

type repository struct {
	cache cache.RedisClient
}

func NewRepository(cache cache.RedisClient) *repository {
	return &repository{
		cache: cache,
	}
}

func (r *repository) getUserSessionCacheKey(sessionUuid *string) string {
	return fmt.Sprintf("%s%s", cacheKeyUserSessionPrefix, *sessionUuid)
}

func (r *repository) getUserSetCacheKey(userUuid *string) string {
	return fmt.Sprintf("%s%s", cacheKeyUserSetPrefix, *userUuid)
}
