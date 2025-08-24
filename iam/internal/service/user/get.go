package user

import (
	"context"

	"github.com/dfg007star/go_rocket/iam/internal/model"
)

func (s *service) Get(ctx context.Context, userUuid *string) (*model.User, error) {
	user, err := s.userRepository.GetByUserUuid(ctx, userUuid)
	if err != nil {
		return nil, model.ErrUserNotFound
	}

	return user, nil
}
