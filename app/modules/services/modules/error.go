package modules

import (
	"github.com/peyman-abdi/bahman/app/interfaces/services"
	"github.com/peyman-abdi/bahman/app/modules/services/errors"
	"reflect"
)

type moduleErrorImpl struct {
	services.ApplicationError
	Message string
	Target  services.Module
}
var _ error = (*moduleErrorImpl)(nil)
var _ services.ApplicationError = (*moduleErrorImpl)(nil)

func (m *moduleErrorImpl) Module() services.Module {
	return m.Target
}

func (m *moduleErrorImpl) StackTrace() string {
	return m.StackTrace()
}

func (m *moduleErrorImpl) File() string {
	return m.File()
}

func (m *moduleErrorImpl) Line() int {
	return m.Line()
}

func (m *moduleErrorImpl) Type() string {
	return m.Type()
}

func (m *moduleErrorImpl) Code() int {
	return m.Code()
}

func (m *moduleErrorImpl) IsErrorType(t interface{}) bool {
	return m.IsErrorType(t)
}

func (m *moduleErrorImpl) Error() string {
	return "Error in module: " + reflect.TypeOf(m.Target).String() + "\n" + m.Message
}

func ModuleError(module services.Module, message string) *moduleErrorImpl {
	return &moduleErrorImpl{
		ApplicationError: errors.NewApplicationError(&moduleManagerImpl{}, 0, message),
		Message: message,
		Target:  module,
	}
}
