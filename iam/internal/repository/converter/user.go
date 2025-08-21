package converter

import (
	"github.com/dfg007star/go_rocket/iam/internal/model"
	repoModel "github.com/dfg007star/go_rocket/iam/internal/repository/model"
)

func RepoModelToUser(user *repoModel.User) *model.User {
	return &model.User{
		UserUuid:            user.UserUuid,
		Login:               user.Login,
		Email:               user.Email,
		Password:            user.Password,
		NotificationMethods: []model.NotificationMethod{},
		CreatedAt:           user.CreatedAt,
		UpdatedAt:           user.UpdatedAt,
	}
}

func UserToRepoModel(user *model.User) *repoModel.User {
	return &repoModel.User{
		UserUuid:            user.UserUuid,
		Login:               user.Login,
		Email:               user.Email,
		Password:            user.Password,
		NotificationMethods: []model.NotificationMethod{},
		CreatedAt:           user.CreatedAt,
		UpdatedAt:           user.UpdatedAt,
	}
}
