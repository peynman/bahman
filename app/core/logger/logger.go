package logger

import (
	"github.com/sirupsen/logrus"
	"github.com/peyman-abdi/avalanche/app/interfaces"
)

type loggerImpl struct {
	loggers []*logrus.Logger
	channels map[string]interfaces.LoggingChannel
}

func Initialize() interfaces.Logger  {
	log := new(loggerImpl)
	log.channels = make(map[string]interfaces.LoggingChannel)
	return log
}

func (l *loggerImpl) LoadChannels(services interfaces.Services) {
	app := services.App()
	config := services.Config()
	/* load channel app */
	modules := app.InitAvalanchePlugins(app.ModulesPath("channels"), services)
	for _, module := range modules {
		logChannel := module.Interface().(interfaces.LoggingChannel)
		l.channels[logChannel.GetChannelName()] = logChannel
	}

	/* setup channel drivers */
	logDrivers := config.GetStringArray("logging.hooks", []string{"console"})
	for _, driverName := range logDrivers {
		driver := l.channels[driverName]
		if driver == nil {
			panic("Driver with name " + driver.GetChannelName() + " not found in log channels")
		}
		if driver.Config("logging.channels." + driver.GetChannelName()) {
			l.loggers = append(l.loggers, driver.GetLogger())
		}
	}
}

func (l *loggerImpl) log(level string, message string, fields map[string]interface{})  {
	switch level {
	case "debug":
		for _, logger := range l.loggers  {
			logger.WithFields(fields).Debug(message)
		}
	case "info":
		for _, logger := range l.loggers  {
			logger.WithFields(fields).Info(message)
		}
	case "warn":
		for _, logger := range l.loggers  {
			logger.WithFields(fields).Warn(message)
		}
	case "error":
		for _, logger := range l.loggers  {
			logger.WithFields(fields).Error(message)
		}
	case "fatal":
		for _, logger := range l.loggers  {
			logger.WithFields(fields).Fatal(message)
		}
	case "panic":
		for _, logger := range l.loggers  {
			logger.WithFields(fields).Panic(message)
		}
	}
}

func (l *loggerImpl) Debug(message string) {
	l.log("debug", message, nil)
}

func (l *loggerImpl) Info(message string) {
	l.log("info", message, nil)
}

func (l *loggerImpl) Warn(message string) {
	l.log("warn", message, nil)
}

func (l *loggerImpl) Error(message string) {
	l.log("error", message, nil)
}

func (l *loggerImpl) Fatal(message string) {
	l.log("fatal", message, nil)
}

func (l *loggerImpl) Panic(message string) {
	l.log("panic", message, nil)
}

func (l *loggerImpl) DebugFields(message string, fields map[string]interface{}) {
	l.log("debug", message, fields)
}

func (l *loggerImpl) InfoFields(message string, fields map[string]interface{}) {
	l.log("info", message, fields)
}

func (l *loggerImpl) WarnFields(message string, fields map[string]interface{}) {
	l.log("warn", message, fields)
}

func (l *loggerImpl) ErrorFields(message string, fields map[string]interface{}) {
	l.log("error", message, fields)
}

func (l *loggerImpl) FatalFields(message string, fields map[string]interface{}) {
	l.log("fatal", message, fields)
}

func (l *loggerImpl) PanicFields(message string, fields map[string]interface{}) {
	l.log("panic", message, fields)
}
