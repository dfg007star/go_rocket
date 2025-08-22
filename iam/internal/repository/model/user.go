package model

import (
	"time"
)

type NotificationMethod struct {
	ProviderName string
	Target       string
}

type User struct {
	UserUuid            *string
	Login               string
	Email               string
	Password            string
	NotificationMethods []NotificationMethod
	CreatedAt           time.Time
	UpdatedAt           *time.Time
}

type UserRedisView struct {
	UserUuid            string `redis:"user_uuid"`
	Login               string `redis:"login"`
	Email               string `redis:"email"`
	Password            string `redis:"password"`
	NotificationMethods string `redis:"notification_methods"`
	CreatedAtNs         int64  `redis:"created_at"`
	UpdatedAtNs         *int64 `redis:"updated_at,omitempty"`
}
