package modules

import (
	"reflect"
	"github.com/peyman-abdi/avalanche/app/interfaces/core"
)

type ModuleErrorImpl struct {
	Message string
	Target  core.Module
}

var _ error = (*ModuleErrorImpl)(nil)

func (m *ModuleErrorImpl) Error() string {
	return "Error in module: " + reflect.TypeOf(m.Target).String() + "\n" + m.Message
}

func ModuleError(module core.Module, message string) *ModuleErrorImpl {
	return &ModuleErrorImpl{
		Message: message,
		Target:  module,
	}
}
