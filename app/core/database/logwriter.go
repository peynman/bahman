package database

import (
	"avalanche/app/core/logger"
	"github.com/sirupsen/logrus"
)

type AvalancheDBLogWriter struct {
}

func (_ *AvalancheDBLogWriter) Println(v... interface{})  {
	if len(v) > 0 {
		logger.ErrorFields(v[0].(string), logrus.Fields{
			"source": "Gorm Database",
		})
	}
}