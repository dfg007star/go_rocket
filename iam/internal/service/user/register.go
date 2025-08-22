package user

import (
	"context"

	"github.com/dfg007star/go_rocket/iam/internal/model"
)

func (s *service) Register(ctx context.Context, user *model.User) (*string, error) {
	userUuid, err := s.userRepository.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return userUuid, nil
}
