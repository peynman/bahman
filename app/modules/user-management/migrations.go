package user_management

import (
	"avalanche/app/core/interfaces"
)

type UsersTable struct {
}
var _ interfaces.Migratable = (*UsersTable)(nil)

func (_ *UsersTable) Up(db interfaces.Migrator) bool {

	return true
}
func (_ *UsersTable) Down(db interfaces.Migrator) bool {
	return true
}

