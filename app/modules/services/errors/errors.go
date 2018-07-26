package errors

import (
	"github.com/peyman-abdi/bahman/app/interfaces/services"
	"runtime"
	"reflect"
)

type applicationErrorImpl struct {
	message string
	typeName string
	stackTrace string
	file string
	line int
	code int
}

var _ services.ApplicationError = (*applicationErrorImpl)(nil)

func (a *applicationErrorImpl) IsErrorType(t interface{}) bool {
	switch t.(type) {
	case string:
		return a.typeName == t
	default:
		return a.typeName == reflect.TypeOf(t).Name()
	}
}

func (a *applicationErrorImpl) File() string {
	return a.file
}

func (a *applicationErrorImpl) Code() int {
	return a.code
}

func (a *applicationErrorImpl) Line() int {
	return a.line
}

func (a *applicationErrorImpl) Type() string {
	return a.typeName
}

func (a *applicationErrorImpl) Error() string {
	return a.message
}

func (a *applicationErrorImpl) StackTrace() string {
	return a.stackTrace
}

func NewApplicationError(t interface{}, code int, m string) services.ApplicationError {
	_, file, line, ok := runtime.Caller(1)
	buf := make([]byte, 1<<16)
	stackSize := runtime.Stack(buf, true)

	if ok {
		return &applicationErrorImpl{
			message: m,
			typeName: reflect.TypeOf(t).Name(),
			code: code,
			stackTrace: string(buf[0:stackSize]),
			file: file,
			line: line,
		}
	}

	return &applicationErrorImpl{
		message: m,
		typeName: reflect.TypeOf(t).Name(),
		code: code,
	}
}

