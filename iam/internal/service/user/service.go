package user

import (
	"github.com/dfg007star/go_rocket/iam/internal/repository/user"
	def "github.com/dfg007star/go_rocket/iam/internal/service"
)

var _ def.UserService = (*service)(nil)

type service struct {
	userRepository repository.UserRepository
}

func NewUserService(
	userRepository repository.UserRepository,
) *service {
	return &service{
		userRepository: userRepository,
	}
}
