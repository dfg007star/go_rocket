package converter

import (
	"encoding/json"
	"time"

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

func UserToRedisView(user *model.User) (*repoModel.UserRedisView, error) {
	methodsJSON, err := json.Marshal(user.NotificationMethods)
	if err != nil {
		return nil, err
	}

	var updatedAtNs *int64
	if user.UpdatedAt != nil {
		ns := user.UpdatedAt.UnixNano()
		updatedAtNs = &ns
	}

	return &repoModel.UserRedisView{
		UserUuid:            *user.UserUuid,
		Login:               user.Login,
		Email:               user.Email,
		Password:            user.Password,
		NotificationMethods: string(methodsJSON),
		CreatedAtNs:         user.CreatedAt.UnixNano(),
		UpdatedAtNs:         updatedAtNs,
	}, nil
}

func RedisViewToUser(redisView *repoModel.UserRedisView) (*model.User, error) {
	var methods []model.NotificationMethod
	if redisView.NotificationMethods != "" {
		err := json.Unmarshal([]byte(redisView.NotificationMethods), &methods)
		if err != nil {
			return nil, err
		}
	}

	createdAt := time.Unix(0, redisView.CreatedAtNs)

	var updatedAt *time.Time
	if redisView.UpdatedAtNs != nil {
		t := time.Unix(0, *redisView.UpdatedAtNs)
		updatedAt = &t
	}

	return &model.User{
		UserUuid:            &redisView.UserUuid,
		Login:               redisView.Login,
		Email:               redisView.Email,
		Password:            redisView.Password,
		NotificationMethods: methods,
		CreatedAt:           createdAt,
		UpdatedAt:           updatedAt,
	}, nil
}
