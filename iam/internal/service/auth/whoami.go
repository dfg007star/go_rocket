package auth

import (
	"context"

	"github.com/dfg007star/go_rocket/iam/internal/model"
)

func (s *service) WhoAmI(ctx context.Context, sessionUuid *string) (*model.User, error) {
	user, err := s.sessionRepository.Get(ctx, &model.Session{SessionUuid: *sessionUuid})
	if err != nil {
		return nil, err
	}

	return user, nil
}
