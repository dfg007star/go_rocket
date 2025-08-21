package auth

import (
	"github.com/dfg007star/go_rocket/iam/internal/repository"
	def "github.com/dfg007star/go_rocket/iam/internal/service"
)

var _ def.AuthService = (*service)(nil)

type service struct {
	sessionRepository repository.SessionRepository
	userRepository    repository.UserRepository
}

func NewAuthService(
	sessionRepository repository.SessionRepository,
	userRepository repository.UserRepository,
) *service {
	return &service{
		sessionRepository: sessionRepository,
		userRepository:    userRepository,
	}
}
