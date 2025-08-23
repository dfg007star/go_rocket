package converter

import (
	"github.com/dfg007star/go_rocket/iam/internal/model"
	authV1 "github.com/dfg007star/go_rocket/shared/pkg/proto/auth/v1"
	commonV1 "github.com/dfg007star/go_rocket/shared/pkg/proto/common/v1"
)

func UserToWhoAmIResponse(user *model.User) *authV1.WhoAmIResponse {
	return &authV1.WhoAmIResponse{
		UserInfo: &commonV1.UserInfo{
			UserUuid: *user.UserUuid,
			Login:    user.Login,
			Email:    user.Email,
		},
	}
}
