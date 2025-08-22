package user

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/dfg007star/go_rocket/iam/internal/model"
	"github.com/dfg007star/go_rocket/iam/internal/repository/converter"
)

func (r *repository) Create(ctx context.Context, user *model.User) (*string, error) {
	repoUser := converter.UserToRepoModel(user)
	notificationMethodsJSON, err := json.Marshal(repoUser.NotificationMethods)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal notification methods: %w", err)
	}

	var userUuid string
	err = r.data.QueryRow(ctx, `
		INSERT INTO users (
			login, 
			email, 
			password, 
			notification_methods,
			created_at,
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
    RETURNING user_uuid`,
		repoUser.Login,
		repoUser.Email,
		repoUser.Password,
		notificationMethodsJSON,
		repoUser.CreatedAt,
	).Scan(&userUuid)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &userUuid, nil
}
