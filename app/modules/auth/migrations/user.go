package migrations

import (
	"github.com/peyman-abdi/bahman/app/interfaces/services"
	"github.com/peyman-abdi/bahman/app/modules/auth/models"
)

type AuthUserMigrate struct {
}

func (*AuthUserMigrate) Up(db services.Migratory) error {
	if err := db.AutoMigrate(&models.UserAuth{}); err != nil {
		return err
	}
	return nil
}

func (*AuthUserMigrate) Down(db services.Migratory) error {
	if err := db.DropTableIfExists(&models.UserAuth{}); err != nil {
		return err
	}
	return nil
}
