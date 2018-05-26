package database

import (
	"github.com/jinzhu/gorm"
	"reflect"
	"github.com/peyman-abdi/avalanche/app/interfaces"
	"github.com/peyman-abdi/avalanche/app/core/modules"
	"strings"
)

type MigrationManager struct {
	migrationsTableName string
	connection *gorm.DB
	log interfaces.Logger
}

func (m *MigrationManager) AutoMigrate(entity interfaces.Entity) error {
 	err := m.connection.AutoMigrate(entity).Error
 	if err != nil {
 		m.log.ErrorFields("Migrator:AutoMigrate", map[string]interface{} {
 			"error": err,
 			"table": entity.TableName(),
		})
	}
	return err
}

func (m *MigrationManager) CreateTable(entity interfaces.Entity) error {
	err := m.connection.CreateTable(entity).Error
	if err != nil {
		m.log.ErrorFields("Migrator:CreateTable", map[string]interface{} {
			"error": err,
			"table": entity.TableName(),
		})
	}
	return err
}

func (m *MigrationManager) DropTable(table... interface{}) error {
	err := m.connection.DropTable(table).Error
	if err != nil {
		m.log.ErrorFields("Migrator:DropTable", map[string]interface{} {
			"error": err,
			"table": table,
		})
	}
	return nil
}

func (m *MigrationManager) DropTableIfExists(table... interface{}) error {
	err := m.connection.DropTableIfExists(table).Error
	if err != nil {
		m.log.ErrorFields("Migrator:DropTableIfExists", map[string]interface{} {
			"error": err,
			"table": table,
		})
	}
	return nil
}

func (m *MigrationManager) HasTable(table interface{}) bool {
	return m.connection.HasTable(table)
}

func (m *MigrationManager) DropColumn(column string) error {
	err := m.connection.DropColumn(column).Error
	if err != nil {
		m.log.ErrorFields("Migrator:DropColumn", map[string]interface{} {
			"error": err,
			"table": column,
		})
	}
	return nil
}

func (m *MigrationManager) Migrate(migrates []interfaces.Migratable) error {
	return m.migrate(m.connection, migrates)
}

func (m *MigrationManager) Rollback(migrates []interfaces.Migratable) error {
	return m.rollback(m.connection, migrates)
}

func (m *MigrationManager) Connection(connection string) interfaces.Migrator {
	conn := connections[connection]
	if conn == nil {
		conn = runtimeConnection
	}

	n := &MigrationManager {
		migrationsTableName: m.migrationsTableName,
		connection: connections[connection],
		log: m.log,
	}
	return n
}

func (m *MigrationManager) setup(connection *gorm.DB) {
	if !connection.HasTable(&MigrationModel{}) {
		connection.AutoMigrate(&MigrationModel{})
	}
	if !connection.HasTable(&modules.ModuleModel{}) {
		connection.AutoMigrate(&modules.ModuleModel{})
	}
}

func (m *MigrationManager) migrate(connection *gorm.DB, migrates []interfaces.Migratable) error {
	var migrations []*MigrationModel
	var migratableInterfaces []string

	for _, migratable := range migrates {
		migratableInterfaces = append(migratableInterfaces, reflect.TypeOf(migratable).String())
	}

	repo := repoManager.Query(&MigrationModel{})
	err := repo.Where("interface IN (`?`)", strings.Join(migratableInterfaces, "`,`")).Get(&migrations)
	if err != nil {
		m.log.ErrorFields("Migrator:Migrate", map[string]interface{} {
			"error": err,
			"connection": reflect.TypeOf(connection).String(),
			"migrates": reflect.TypeOf(migrates).String(),
		})
		return err
	}

	var maxStep int
	err = repo.Select("max(step)").Get(&maxStep)
	if err != nil {
		m.log.ErrorFields("Migrator:Migrate", map[string]interface{} {
			"error": err,
			"connection": reflect.TypeOf(connection).String(),
			"migrates": reflect.TypeOf(migrates).String(),
		})
		return err
	}

	for index, migratable := range migrates {
		if migration := getMigration(migrations, migratable); migration == nil {
			if migratable.Up(m) {
				migration = &MigrationModel{
					Step: maxStep + 1,
					Interface: migratableInterfaces[index],
				}
				err = repoManager.Insert(migration)
				if err != nil {
					m.log.ErrorFields("Migrator:Migrate:Insert", map[string]interface{} {
						"error": err,
						"connection": reflect.TypeOf(connection).String(),
						"migrates": reflect.TypeOf(migrates).String(),
					})
					return err
				}
			}
		}
	}

	return nil
}

func (m *MigrationManager) rollback(connection *gorm.DB, migrates []interfaces.Migratable) error {
	var migrations []*MigrationModel
	var migratableInterfaces []string

	for _, migratable := range migrates {
		migratableInterfaces = append(migratableInterfaces, reflect.TypeOf(migratable).String())
	}

	repo := repoManager.Query(&MigrationModel{})

	err := repo.Where("interface IN (?)", migratableInterfaces).Get(&migrations)
	if err != nil {
		m.log.ErrorFields("Migrator:Rollback", map[string]interface{} {
			"error": err,
			"connection": reflect.TypeOf(connection).String(),
			"migrates": reflect.TypeOf(migrates).String(),
		})
		return err
	}

	var maxStep int
	err = repo.Select("max(step)").GetValue(&maxStep)
	if err != nil {
		m.log.ErrorFields("Migrator:Rollback", map[string]interface{} {
			"error": err,
			"connection": reflect.TypeOf(connection).String(),
			"migrates": reflect.TypeOf(migrates).String(),
		})
		return err
	}

	for _, migratable := range migrates {
		if migration := getMigration(migrations, migratable); migration == nil {
			if migratable.Down(m) {
				err = repoManager.DeleteEntity(migration)
				if err != nil {
					m.log.ErrorFields("Migrator:Rollback", map[string]interface{} {
						"error": err,
						"connection": reflect.TypeOf(connection).String(),
						"migrates": reflect.TypeOf(migrates).String(),
					})
					return err
				}
			}
		}
	}

	return nil
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



