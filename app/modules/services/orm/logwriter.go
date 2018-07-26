package orm

import (
	"fmt"
	"github.com/peyman-abdi/bahman/app/interfaces/services"
)

type bahmanDBLogWriter struct {
	logger services.Logger
}

func (a *bahmanDBLogWriter) Println(v ...interface{}) {
	if a.logger != nil {
		var data = make(map[string]interface{})

		if len(v) > 0 {
			for index, d := range v {
				data[fmt.Sprintf("%d", index)] = fmt.Sprintf("%v", d)
			}
		}

		a.logger.DebugFields("Database Query", data)
	}
}

func NewLogWriter(logger services.Logger) *bahmanDBLogWriter {
	a := new(bahmanDBLogWriter)
	a.logger = logger
	return a
}