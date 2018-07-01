package database

import (
	"github.com/peyman-abdi/avalanche/app/interfaces/services"
	"time"
)

type MigrationModel struct {
	ID        int64
	Step      int64
	Interface string `gorm:"size:192,unique_index"`
	CreatedAt *time.Time
}

var _ services.Entity = (*MigrationModel)(nil)
var _ services.EntityConnection = (*MigrationModel)(nil)

func (m *MigrationModel) TableName() string {
	return "migrations"
}
func (m *MigrationModel) PrimaryKey() string {
	return "id"
}
func (m *MigrationModel) ConnectionName() string {
	return "runtime"
}
