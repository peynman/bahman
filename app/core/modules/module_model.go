package modules

import (
	"time"
	"github.com/peyman-abdi/avalanche/app/interfaces"
)

const (
	INSTALLED = 1 << iota
	ACTIVE = 1 << iota
)

type ModuleModel struct {
	ID int64
	Flags int
	Interface string `gorm:"size:192,unique_index"`
	CreatedAt *time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}
var _ interfaces.Entity = (*ModuleModel)(nil)
var _ interfaces.EntityConnection = (*ModuleModel)(nil)

func (m *ModuleModel) TableName() string {
	return "modules"
}
func (m *ModuleModel) PrimaryKey() string {
	return "id"
}
func (m *ModuleModel) ConnectionName() string {
	return "runtime"
}