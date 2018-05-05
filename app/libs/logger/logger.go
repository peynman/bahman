package logger

import (
	"github.com/sirupsen/logrus"
	"os"
	"avalanche/app/libs"
	"path/filepath"
	plugin2 "plugin"
	"strings"
	"avalanche/app/libs/interfaces"
	"AvConfig/lib"
)

var (
	loggers []logrus.Logger
	channels = make(map[string]interfaces.LoggingChannel)
)

func Initialize()  {
	/* load channel plugins */
	var moduleFiles []string
	channelModulesDir := app.ModulesPath("channels")
	filepath.Walk(channelModulesDir, func(path string, info os.FileInfo, err error) error {
		moduleFiles = append(moduleFiles, path)
		return nil
	})
	for _, moduleFile := range moduleFiles {
		if !strings.HasSuffix(moduleFile, ".so") {
			continue
		}

		plugin, err := plugin2.Open(moduleFile)
		if err != nil {
			panic(err)
		}

		symInstance, err := plugin.Lookup("GetInstance")
		if err != nil {
			panic(err)
		}

		instance := symInstance.(func() interfaces.LoggingChannel)
		filename := filepath.Base(moduleFile)
		extension := filepath.Ext(moduleFile)
		driverName := filename[0:len(filename)-len(extension)]
		channels[driverName] = instance()
	}

	/* setup channel drivers */
	logDrivers := config.GetStringArray("logging.hooks", []string{"console"})
	for _, driverName := range logDrivers {
		driver := channels[driverName]
		if driver == nil {
			panic("Driver with name " + driverName + " not found in log channels")
		}
		if driver.Config("logging.channels." + driverName) {
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
