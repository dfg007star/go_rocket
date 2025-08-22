package converter

import (
	"github.com/dfg007star/go_rocket/iam/internal/model"
	commonV1 "github.com/dfg007star/go_rocket/shared/pkg/proto/common/v1"
	userV1 "github.com/dfg007star/go_rocket/shared/pkg/proto/user/v1"
)

func UserToGetUserResponse(user *model.User) *userV1.GetUserResponse {
	notificationMethods := make([]*commonV1.NotificationMethod, 0, len(user.NotificationMethods))
	for _, method := range user.NotificationMethods {
		notificationMethods = append(notificationMethods, &commonV1.NotificationMethod{
			ProviderName: method.ProviderName,
			Target:       method.Target,
		})
	}

	return &userV1.GetUserResponse{
		UserInfo: &commonV1.UserInfo{
			UserUuid: *user.UserUuid,
			Login:    user.Login,
			Email:    user.Email,
		},
		NotificationMethods: notificationMethods,
	}
}
