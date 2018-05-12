package logger

import (
	"github.com/sirupsen/logrus"
	"avalanche/app/core/interfaces"
	"avalanche/app/core/app"
	"avalanche/app/core/config"
)

var (
	loggers []*logrus.Logger
	channels = make(map[string]interfaces.LoggingChannel)
)

func Initialize()  {
	/* load channel app */
	modules := app.InitAvalanchePlugins(app.ModulesPath("channels"))
	for _, module := range modules {
		logChannel := module.Interface().(interfaces.LoggingChannel)
		channels[logChannel.GetChannelName()] = logChannel
	}

	/* setup channel drivers */
	logDrivers := config.GetStringArray("logging.hooks", []string{"console"})
	for _, driverName := range logDrivers {
		driver := channels[driverName]
		if driver == nil {
			panic("Driver with name " + driver.GetChannelName() + " not found in log channels")
		}
		if driver.Config("logging.channels." + driver.GetChannelName()) {
			loggers = append(loggers, driver.GetLogger())
		}
	}
}

func log(level string, message string, fields logrus.Fields)  {
	switch level {
	case "debug":
		for _, logger := range loggers  {
			logger.WithFields(fields).Debug(message)
		}
	case "info":
		for _, logger := range loggers  {
			logger.WithFields(fields).Info(message)
		}
	case "warn":
		for _, logger := range loggers  {
			logger.WithFields(fields).Warn(message)
		}
	case "error":
		for _, logger := range loggers  {
			logger.WithFields(fields).Error(message)
		}
	case "fatal":
		for _, logger := range loggers  {
			logger.WithFields(fields).Fatal(message)
		}
	case "panic":
		for _, logger := range loggers  {
			logger.WithFields(fields).Panic(message)
		}
	}
}

func Debug(message string) {
	log("debug", message, nil)
}

func Info(message string) {
	log("info", message, nil)
}

func Warn(message string) {
	log("warn", message, nil)
}

func Error(message string) {
	log("error", message, nil)
}

func Fatal(message string) {
	log("fatal", message, nil)
}

func Panic(message string) {
	log("panic", message, nil)
}

func DebugFields(message string, fields logrus.Fields) {
	log("debug", message, fields)
}

func InfoFields(message string, fields logrus.Fields) {
	log("info", message, fields)
}

func WarnFields(message string, fields logrus.Fields) {
	log("warn", message, fields)
}

func ErrorFields(message string, fields logrus.Fields) {
	log("error", message, fields)
}

func FatalFields(message string, fields logrus.Fields) {
	log("fatal", message, fields)
}

func PanicFields(message string, fields logrus.Fields) {
	log("panic", message, fields)
}
