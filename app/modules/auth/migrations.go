package auth

import "github.com/peyman-abdi/avalanche/app/interfaces/services"

type MigrateModels struct {
}

func (*MigrateModels) Up(db services.Migrator) error {
	if err := db.AutoMigrate(&PermissionModel{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&RoleModel{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&UserModel{}); err != nil {
		return err
	}
	return nil
}

func (*MigrateModels) Down(db services.Migrator) error {
	if err := db.DropTableIfExists(&UserModel{}); err != nil {
		return err
	}
	if err := db.DropTableIfExists(&RoleModel{}); err != nil {
		return err
	}
	if err := db.DropTableIfExists(&PermissionModel{}); err != nil {
		return err
	}
	return nil
}

func (*AuthenticationModule) Migrations() []services.Migratable {
	return []services.Migratable{
		new(MigrateModels),
	}
}
