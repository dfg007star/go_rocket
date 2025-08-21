package session

import (
	"context"
	"fmt"

	"github.com/dfg007star/go_rocket/iam/internal/model"
)

const (
	cacheKeyUserSessionPrefix = "user:session:"
	cacheKeyUserSetPrefix     = "user:set:"
)

func (r *repository) getCacheKey(sessionUuid string) string {
	return fmt.Sprintf("%s%s", cacheKeyUserSessionPrefix, sessionUuid)
}

func (r *repository) Get(ctx context.Context, session *model.Session) (*model.User, error) {
	cacheKey := r.getCacheKey(session.SessionUuid)

	values, err := r.cache.HGetAll(ctx, cacheKey)
}
