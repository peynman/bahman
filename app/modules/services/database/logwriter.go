package database

import "fmt"

type AvalancheDBLogWriter struct {
}

func (_ *AvalancheDBLogWriter) Println(v ...interface{}) {
	if loggerRef != nil {
		var data = make(map[string]interface{})

		if len(v) > 0 {
			for index, d := range v {
				data[fmt.Sprintf("%d", index)] = fmt.Sprintf("%v", d)
			}
		}

		loggerRef.DebugFields("Database Query", data)
	}
}
