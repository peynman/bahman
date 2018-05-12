package database

import (
	"time"
	"github.com/jinzhu/gorm"
	"reflect"
	"avalanche/app/core/interfaces"
	"strings"
)

type Migration struct {
	ID int64
	Step int
	Interface string `gorm:"size:192,unique_index"`
	CreatedAt *time.Time
}
type MigrationManager struct {
	tableName string
}

func (m *MigrationManager) Setup(connection *gorm.DB) {
	if !connection.HasTable(&Migration{}) {
		connection.AutoMigrate(&Migration{})
	}
}

func (m *MigrationManager) Migrate(connection *gorm.DB, migratables []interfaces.Migratable) {
	var migrations []*Migration
	var migratableInterfaces []string

	for _, migratable := range migratables {
		migratableInterfaces = append(migratableInterfaces, reflect.TypeOf(migratable).String())
	}

	query := "interface IN (`" + strings.Join(migratableInterfaces, "`,`") + "`"
	connection.Find(&migrations).Where(query)

	row := connection.Raw("SELECT max(step) FROM migrations").Row()
	var maxStep int
 	row.Scan(&maxStep)

	for index, migratable := range migratables {
		if migration := getMigration(migrations, migratable); migration == nil {
			if migratable.Up(m) {
				migration = &Migration{
					Step: maxStep + 1,
					Interface: migratableInterfaces[index],
				}
				connection.Create(migration)
			}
		}
	}
}

func (m *MigrationManager) Rollback(connection *gorm.DB, migratables []interfaces.Migratable) {
	var migrations []*Migration
	var migratableInterfaces []string

	for _, migratable := range migratables {
		migratableInterfaces = append(migratableInterfaces, reflect.TypeOf(migratable).String())
	}

	query := "interface IN (`" + strings.Join(migratableInterfaces, "`,`") + "`"
	connection.Find(&migrations).Where(query)

	row := connection.Raw("SELECT max(step) FROM migrations").Row()
	var maxStep int
	row.Scan(&maxStep)

	for _, migratable := range migratables {
		if migration := getMigration(migrations, migratable); migration == nil {
			if migratable.Down(m) {
				connection.Delete(migration)
			}
		}
	}
}

func (m *MigrationManager) Reset(connection *gorm.DB) {

}

func getMigration(migrations []*Migration, migratable interfaces.Migratable) *Migration {
	interfaceName := reflect.TypeOf(migratable).String()
	for _, migration := range migrations {
		if migration.Interface == interfaceName {
			return migration
		}
	}
	return nil
}

