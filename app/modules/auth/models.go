package auth

import "time"

// UserModel
type UserModel struct {
	ID        uint64
	Flags     int
	CreatedAt *time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
	Roles     []*RoleModel `gorm:"many2many:user_roles"`
}
func (*UserModel) TableName() string { return "users" }

// RoleModel
type RoleModel struct {
	ID          uint64
	Title       string
	Description *string
	Hash        string `gorm:"UNIQUE_INDEX"`
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
	DeletedAt   *time.Time
	Permissions []*PermissionModel `gorm:"many2many:role_permissions"`
}
func (*RoleModel) TableName() string { return "roles" }

// PermissionModel
type PermissionModel struct {
	ID          uint64
	Title       string
	Description *string
	Hash        string  `gorm:"UNIQUE_INDEX"`
}
func (*PermissionModel) TableName() string { return "permissions" }
