package database

import (
	"time"
	"github.com/jinzhu/gorm"
	"reflect"
	"strings"
	"fmt"
	"avalanche/app/core/interfaces"
)

type Migration struct {
	ID int64
	Step int
	Interface string `gorm:"size:192,unique"`
	CreatedAt *time.Time
}
type MigrationManager struct {
}

func (m *MigrationManager) Setup(connection *gorm.DB) {
	if !connection.HasTable(&Migration{}) {
		connection.AutoMigrate(&Migration{})
	}
}

func (m *MigrationManager) Migrate(connection *gorm.DB, migratables []interfaces.Migratable) {
	var migrations []Migration
	var migratableInterfaces []string

	for _, migratable := range migratables {
		migratableInterfaces = append(migratableInterfaces, reflect.TypeOf(migratable).String())
	}

	query := "interface IN (`" + strings.Join(migratableInterfaces, "`,`") + "`"
	fmt.Println(query)
	if err := connection.Find(&migrations, query).Error; err != nil {
		fmt.Println("DB Error")
		fmt.Println(err)
	}

	//if migratable.Up(connection) {
	//	migration := Migration {
	//		Step: 1,
	//
	//	}
	//}
}

func (m *MigrationManager) Rollback(connection *gorm.DB, steps int) {

}

func (m *MigrationManager) Reset(connection *gorm.DB) {

}


