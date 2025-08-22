package auth

import (
	"context"

	"github.com/dfg007star/go_rocket/iam/internal/config"
	"github.com/dfg007star/go_rocket/iam/internal/model"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (s *service) Login(ctx context.Context, login *string, password *string) (*string, error) {
	user, err := s.userRepository.GetByUserLogin(ctx, login)
	if err != nil {
		return nil, model.ErrUserNotFound
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(*password))
	if err != nil {
		return nil, model.ErrUserPassword
	}

	sessionUuid := uuid.New().String()
	err = s.sessionRepository.Create(ctx, &sessionUuid, user, config.AppConfig().Session.SessionTtl())
	if err != nil {
		return nil, model.ErrSessionCreate
	}

	return &sessionUuid, nil
}
