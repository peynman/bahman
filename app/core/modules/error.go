package modules

import (
	"reflect"
	"github.com/peyman-abdi/avalanche/app/interfaces"
)

type ModuleErrorImpl struct {
	Message string
	Target interfaces.Module
}
var _ error = (*ModuleErrorImpl)(nil)

func (m *ModuleErrorImpl) Error() string {
	return "Error in module: " + reflect.TypeOf(m.Target).String() + "\n" + m.Message
}

func ModuleError(module interfaces.Module, message string) *ModuleErrorImpl {
	return &ModuleErrorImpl{
		Message: message,
		Target: module,
	}
}
