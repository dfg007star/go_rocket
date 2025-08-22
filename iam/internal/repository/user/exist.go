package user

import (
	"context"
	"fmt"
)

func (r *repository) CheckLoginExists(ctx context.Context, login *string) (bool, error) {
	var exists bool

	err := r.data.QueryRow(ctx, `
		SELECT EXISTS(
			SELECT 1 FROM users 
			WHERE login = $1
		)`,
		login,
	).Scan(&exists)

	if err != nil {
		return false, fmt.Errorf("failed to check login existence: %w", err)
	}

	return exists, nil
}

func (r *repository) CheckEmailExists(ctx context.Context, email *string) (bool, error) {
	var exists bool

	err := r.data.QueryRow(ctx, `
		SELECT EXISTS(
			SELECT 1 FROM users 
			WHERE email = $1
		)`,
		email,
	).Scan(&exists)

	if err != nil {
		return false, fmt.Errorf("failed to check email existence: %w", err)
	}

	return exists, nil
}
