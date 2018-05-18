package database

import (
	"github.com/jinzhu/gorm"
	"reflect"
	"github.com/peyman-abdi/avalanche/app/interfaces"
)

type MigrationManager struct {
	migrationsTableName string
	connection *gorm.DB
}

func (m *MigrationManager) AutoMigrate(table... interface{}) error {
	m.connection.AutoMigrate(table)
	return nil
}

func (m *MigrationManager) CreateTable(table... interface{}) error {
	m.connection.CreateTable(table)
	return nil
}

func (m *MigrationManager) DropTable(table... interface{}) error {
	m.connection.DropTable(table)
	return nil
}

func (m *MigrationManager) DropTableIfExists(table... interface{}) error {
	m.connection.DropTableIfExists(table)
	return nil
}

func (m *MigrationManager) DropColumn(column string) error {
	m.connection.DropColumn(column)
	return nil
}

func (m *MigrationManager) Migrate(migrates []interfaces.Migratable) error {
	m.migrate(m.connection, migrates)
	return nil
}

func (m *MigrationManager) Rollback(migrates []interfaces.Migratable) error {
	m.rollback(m.connection, migrates)
	return nil
}

func (m *MigrationManager) Connection(connection string) interfaces.Migrator {
	conn := connections[connection]
	if conn == nil {
		conn = appConnection
	}

	n := &MigrationManager {
		migrationsTableName: m.migrationsTableName,
		connection: connections[connection],
	}
	return n
}

func (m *MigrationManager) setup(connection *gorm.DB) {
	if !connection.HasTable(&MigrationModel{}) {
		connection.AutoMigrate(&MigrationModel{})
	}
}

func (m *MigrationManager) migrate(connection *gorm.DB, migrates []interfaces.Migratable) {
	var migrations []*MigrationModel
	var migratableInterfaces []string

	for _, migratable := range migrates {
		migratableInterfaces = append(migratableInterfaces, reflect.TypeOf(migratable).String())
	}

	repo := repoManager.Query(&MigrationModel{})

	repo.Where("interface IN (?)", migratableInterfaces).Get(&migrations)

	var maxStep int
	repo.Select("max(step)").Get(&maxStep)

	for index, migratable := range migrates {
		if migration := getMigration(migrations, migratable); migration == nil {
			if migratable.Up(m) {
				migration = &MigrationModel{
					Step: maxStep + 1,
					Interface: migratableInterfaces[index],
				}
				connection.Create(migration)
			}
		}
	}
}

func (m *MigrationManager) rollback(connection *gorm.DB, migrates []interfaces.Migratable) {
	var migrations []*MigrationModel
	var migratableInterfaces []string

	for _, migratable := range migrates {
		migratableInterfaces = append(migratableInterfaces, reflect.TypeOf(migratable).String())
	}

	repo := repoManager.Query(&MigrationModel{})

	repo.Where("interface IN (?)", migratableInterfaces).Get(&migrations)

	var maxStep int
	repo.Select("max(step)").Get(&maxStep)

	for _, migratable := range migrates {
		if migration := getMigration(migrations, migratable); migration == nil {
			if migratable.Down(m) {
				connection.Delete(migration)
			}
		}
	}
}

func getMigration(migrations []*MigrationModel, migratable interfaces.Migratable) *MigrationModel {
	interfaceName := reflect.TypeOf(migratable).String()
	for _, migration := range migrations {
		if migration.Interface == interfaceName {
			return migration
		}
	}
	return nil
}



