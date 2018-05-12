package database

import "fmt"

type AvalancheDBLogWriter struct {
}

func (_ *AvalancheDBLogWriter) Println(v... interface{})  {
	if len(v) > 0 {
		fmt.Println(v)
		//logger.ErrorFields("Database Error", logrus.Fields{
		//	"source": "Gorm Database",
		//})
	}
}