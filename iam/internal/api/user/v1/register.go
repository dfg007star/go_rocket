package v1

import (
	"context"
	"time"

	"github.com/dfg007star/go_rocket/iam/internal/model"
	userV1 "github.com/dfg007star/go_rocket/shared/pkg/proto/user/v1"
)

func (a *api) Register(ctx context.Context, req *userV1.RegisterRequest) (*userV1.RegisterResponse, error) {
	notificationMethods := make([]*model.NotificationMethod, 0, len(req.NotificationMethods))
	for _, method := range req.GetNotificationMethods() {
		notificationMethods = append(notificationMethods, &model.NotificationMethod{ProviderName: method.ProviderName, Target: method.Target})
	}
	user := model.User{
		Login:               req.Login,
		Email:               req.Email,
		Password:            req.Password,
		NotificationMethods: notificationMethods,
		CreatedAt:           time.Now(),
	}

	userUuid, err := a.userService.Register(ctx, &user)
	if err != nil {
		return nil, err
	}

	return &userV1.RegisterResponse{UserUuid: *userUuid}, nil
}
