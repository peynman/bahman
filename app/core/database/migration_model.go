package database

import (
	"time"
	"github.com/peyman-abdi/avalanche/app/interfaces"
)

type MigrationModel struct {
	ID int64
	Step int64
	Interface string `gorm:"size:192,unique_index"`
	CreatedAt *time.Time
}
var _ interfaces.Entity = (*MigrationModel)(nil)
var _ interfaces.EntityConnection = (*MigrationModel)(nil)

func (m *MigrationModel) TableName() string {
	return "migrations"
}
func (m *MigrationModel) PrimaryKey() string {
	return "id"
}
func (m *MigrationModel) ConnectionName() string {
	return "runtime"
}