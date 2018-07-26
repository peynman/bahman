package models

import (
	"time"
	"fmt"
)

// UserAuth
type UserAuth struct {
	ID        uint64
	Username  string `gorm:"type:varchar(200);not null;unique"`
	Password  string `gorm:"type:varchar(200);not null"`
	Salt      string `gorm:"type:varchar(200);not null"`
	CreatedAt *time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}

func (a *UserAuth) UID() string {
	return fmt.Sprintf("%d", a.ID)
}

func (*UserAuth) TableName() string { return "users" }

