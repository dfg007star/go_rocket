package session

import (
	"context"
	"time"

	"github.com/dfg007star/go_rocket/iam/internal/model"
	repoConverter "github.com/dfg007star/go_rocket/iam/internal/repository/converter"
)

func (r *repository) Create(ctx context.Context, sessionUuid *string, user *model.User, ttl time.Duration) error {
	sessionKey := r.getUserSessionCacheKey(sessionUuid)

	err := r.cache.Set(ctx, sessionKey, *user.UserUuid)
	if err != nil {
		return err
	}

	redisUser, err := repoConverter.UserToRedisView(user)
	if err != nil {
		return err
	}

	userSetKey := r.getUserSetCacheKey(user.UserUuid)
	err = r.cache.HashSet(ctx, userSetKey, redisUser)
	if err != nil {
		return err
	}

	return nil
}
