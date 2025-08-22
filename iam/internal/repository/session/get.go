package session

import (
	"context"

	repoConverter "github.com/dfg007star/go_rocket/iam/internal/repository/converter"
	repoModel "github.com/dfg007star/go_rocket/iam/internal/repository/model"
	redigo "github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"

	"github.com/dfg007star/go_rocket/iam/internal/model"
)

func (r *repository) Get(ctx context.Context, session *model.Session) (*model.User, error) {
	sessionKey := r.getUserSessionCacheKey(&session.SessionUuid)
	userUuidValue, err := r.cache.Get(ctx, sessionKey)
	if err != nil {
		if errors.Is(err, redigo.ErrNil) {
			return nil, model.ErrSessionNotFound
		}
		return nil, err
	}

	userUuid := string(userUuidValue)
	userKey := r.getUserSetCacheKey(&userUuid)
	redisUserValues, err := r.cache.HGetAll(ctx, userKey)
	if err != nil {
		if errors.Is(err, redigo.ErrNil) {
			return nil, model.ErrUserNotFound
		}
		return nil, err
	}

	if len(redisUserValues) == 0 {
		return nil, model.ErrUserNotFound
	}

	var userRedisView repoModel.UserRedisView
	err = redigo.ScanStruct(redisUserValues, &userRedisView)
	if err != nil {
		return nil, err
	}

	user, err := repoConverter.RedisViewToUser(&userRedisView)
	if err != nil {
		return nil, err
	}

	return user, nil
}
