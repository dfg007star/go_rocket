package user

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/dfg007star/go_rocket/iam/internal/model"
)

func (s *service) Register(ctx context.Context, user *model.User) (*string, error) {
	if err := s.validateUser(user); err != nil {
		return nil, err
	}

	exists, err := s.userRepository.CheckLoginExists(ctx, &user.Login)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, model.ErrUserLoginExists
	}

	exists, err = s.userRepository.CheckEmailExists(ctx, &user.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, model.ErrUserEmailExists
	}

	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user.Password = hashedPassword
	now := time.Now()
	user.CreatedAt = now
	userUuid, err := s.userRepository.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return userUuid, nil
}

func (s *service) validateUser(user *model.User) error {
	if strings.TrimSpace(user.Login) == "" {
		return model.ErrUserLoginRequired
	}
	if len(user.Login) < 3 || len(user.Login) > 50 {
		return model.ErrUserLoginInvalid
	}

	if strings.TrimSpace(user.Email) == "" {
		return model.ErrUserEmailRequired
	}
	if !isValidEmail(user.Email) {
		return model.ErrUserEmailInvalid
	}

	if strings.TrimSpace(user.Password) == "" {
		return model.ErrUserPasswordRequired
	}
	if len(user.Password) < 8 {
		return model.ErrUserPasswordTooShort
	}

	for _, method := range user.NotificationMethods {
		if err := s.validateNotificationMethod(method); err != nil {
			return err
		}
	}

	return nil
}

func (s *service) validateNotificationMethod(method *model.NotificationMethod) error {
	if method == nil {
		return model.ErrNotificationMethodInvalid
	}
	if strings.TrimSpace(method.ProviderName) == "" {
		return model.ErrNotificationMethodInvalid
	}
	if strings.TrimSpace(method.Target) == "" {
		return model.ErrNotificationMethodInvalid
	}
	return nil
}

func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	return emailRegex.MatchString(strings.ToLower(email))
}

func hashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}
