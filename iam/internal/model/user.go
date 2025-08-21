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
