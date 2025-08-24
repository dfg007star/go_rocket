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

func NewSessionRepository(cache cache.RedisClient) *repository {
	return &repository{
		cache: cache,
	}
}

func (r *repository) getUserSessionCacheKey(sessionUuid *string) string {
	if sessionUuid == nil {
		return cacheKeyUserSessionPrefix + "invalid"
	}
	return fmt.Sprintf("%s%s", cacheKeyUserSessionPrefix, *sessionUuid)
}

func (r *repository) getUserSetCacheKey(userUuid *string) string {
	if userUuid == nil {
		return cacheKeyUserSessionPrefix + "invalid"
	}
	return fmt.Sprintf("%s%s", cacheKeyUserSetPrefix, *userUuid)
}
