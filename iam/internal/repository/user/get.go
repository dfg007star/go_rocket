package user

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/dfg007star/go_rocket/iam/internal/model"
	"github.com/dfg007star/go_rocket/iam/internal/repository/converter"
	repoModel "github.com/dfg007star/go_rocket/iam/internal/repository/model"
)

func (r *repository) Get(ctx context.Context, userUuid *string) (*model.User, error) {
	var repoUser repoModel.User
	var notificationMethodsJSON []byte

	err := r.data.QueryRow(ctx, `
		SELECT 
			user_uuid, 
			login, 
			email, 
			password, 
			notification_methods,
			created_at,
			updated_at
		FROM users 
		WHERE user_uuid = $1`,
		userUuid,
	).Scan(
		&repoUser.UserUuid,
		&repoUser.Login,
		&repoUser.Email,
		&repoUser.Password,
		&notificationMethodsJSON,
		&repoUser.CreatedAt,
		&repoUser.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if err := json.Unmarshal(notificationMethodsJSON, &repoUser.NotificationMethods); err != nil {
		return nil, fmt.Errorf("failed to unmarshal notification methods: %w", err)
	}

	return converter.RepoModelToUser(&repoUser), nil
}
