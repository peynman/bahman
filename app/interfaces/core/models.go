package core

type Permission interface {
	Title() string
	Description() string
	Hash() string
}

type Role interface {
	Permissions() []Permission
	HasPermission(permission Permission) bool
	HasPermissions(permission []Permission) bool
}

type User interface {
	Roles() []Role

	HasRole(role Role) bool
	HasPermission(permission Permission) bool
	HasPermissions(permissions []Permission) bool
}

type Authenticates interface {
}

type Authenticator interface {
	FindUserWithCredentials(interface{}) User
	FindUserWithToken(token string) User
	IsCredentialsValid(interface{}) bool
}
