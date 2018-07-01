package modules

import (
	"github.com/peyman-abdi/avalanche/app/interfaces/services"
	"reflect"
)

type ModuleErrorImpl struct {
	Message string
	Target  services.Module
}

var _ error = (*ModuleErrorImpl)(nil)

func (m *ModuleErrorImpl) Error() string {
	return "Error in module: " + reflect.TypeOf(m.Target).String() + "\n" + m.Message
}

func ModuleError(module services.Module, message string) *ModuleErrorImpl {
	return &ModuleErrorImpl{
		Message: message,
		Target:  module,
	}
}
