package converter

import (
	"github.com/dfg007star/go_rocket/iam/internal/model"
	repoModel "github.com/dfg007star/go_rocket/iam/internal/repository/model"
)

func RepoModelToUser(user *repoModel.User) *model.User {
	notificationMethods := make([]model.NotificationMethod, 0, len(user.NotificationMethods))
	for _, method := range user.NotificationMethods {
		notificationMethods = append(notificationMethods, model.NotificationMethod{
			ProviderName: method.ProviderName,
			Target:       method.Target,
		})
	}

	return &model.User{
		UserUuid:            user.UserUuid,
		Login:               user.Login,
		Email:               user.Email,
		Password:            user.Password,
		NotificationMethods: notificationMethods,
		CreatedAt:           user.CreatedAt,
		UpdatedAt:           user.UpdatedAt,
	}
}

func UserToRepoModel(user *model.User) *repoModel.User {
	notificationMethods := make([]repoModel.NotificationMethod, 0, len(user.NotificationMethods))
	for _, method := range user.NotificationMethods {
		notificationMethods = append(notificationMethods, repoModel.NotificationMethod{
			ProviderName: method.ProviderName,
			Target:       method.Target,
		})
	}

	return &repoModel.User{
		UserUuid:            user.UserUuid,
		Login:               user.Login,
		Email:               user.Email,
		Password:            user.Password,
		NotificationMethods: notificationMethods,
		CreatedAt:           user.CreatedAt,
		UpdatedAt:           user.UpdatedAt,
	}
}
