package logger

import (
	"github.com/sirupsen/logrus"
	"runtime"
	"github.com/peyman-abdi/avalanche/app/interfaces/core"
)

type loggerImpl struct {
	loggers         []*logrus.Logger
	channels        map[string]core.LoggingChannel
	appendDebugData func(fields *map[string]interface{})
}

func Initialize(config core.Config) core.Logger {
	log := new(loggerImpl)
	log.channels = make(map[string]core.LoggingChannel)

	if config.GetBoolean("app.debug", false) {
		log.appendDebugData = appendCallerInfo
	} else {
		log.appendDebugData = func(fields *map[string]interface{}) {
		}
	}

	return log
}

func (l *loggerImpl) LoadChannels(services core.Services) {
	app := services.App()
	config := services.Config()
	/* load channel app */
	modules := app.InitAvalanchePlugins(app.ModulesPath("channels"), services)
	for _, module := range modules {
		logChannel := module.Interface().(core.LoggingChannel)
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

func (l *loggerImpl) LoadConsole() {
	l.loggers = append(l.loggers, logrus.New())
}

func (l *loggerImpl) log(level string, message string, fields map[string]interface{}) {
	switch level {
	case "debug":
		for _, logger := range l.loggers {
			logger.WithFields(fields).Debug(message)
		}
	case "info":
		for _, logger := range l.loggers {
			logger.WithFields(fields).Info(message)
		}
	case "warn":
		for _, logger := range l.loggers {
			logger.WithFields(fields).Warn(message)
		}
	case "error":
		for _, logger := range l.loggers {
			logger.WithFields(fields).Error(message)
		}
	case "fatal":
		for _, logger := range l.loggers {
			logger.WithFields(fields).Fatal(message)
		}
	case "panic":
		for _, logger := range l.loggers {
			logger.WithFields(fields).Panic(message)
		}
	}
}

func (l *loggerImpl) Debug(message string) {
	fields := make(map[string]interface{})
	l.appendDebugData(&fields)
	l.log("debug", message, fields)
}

func (l *loggerImpl) Info(message string) {
	fields := make(map[string]interface{})
	l.appendDebugData(&fields)
	l.log("info", message, fields)
}

func (l *loggerImpl) Warn(message string) {
	fields := make(map[string]interface{})
	l.appendDebugData(&fields)
	l.log("warn", message, fields)
}

func (l *loggerImpl) Error(message string) {
	fields := make(map[string]interface{})
	l.appendDebugData(&fields)
	l.log("error", message, fields)
}

func (l *loggerImpl) Fatal(message string) {
	fields := make(map[string]interface{})
	l.appendDebugData(&fields)
	l.log("fatal", message, fields)
}

func (l *loggerImpl) Panic(message string) {
	fields := make(map[string]interface{})
	l.appendDebugData(&fields)
	l.log("panic", message, fields)
}

func (l *loggerImpl) DebugFields(message string, fields map[string]interface{}) {
	l.appendDebugData(&fields)
	l.log("debug", message, fields)
}

func (l *loggerImpl) InfoFields(message string, fields map[string]interface{}) {
	l.appendDebugData(&fields)
	l.log("info", message, fields)
}

func (l *loggerImpl) WarnFields(message string, fields map[string]interface{}) {
	l.appendDebugData(&fields)
	l.log("warn", message, fields)
}

func (l *loggerImpl) ErrorFields(message string, fields map[string]interface{}) {
	l.appendDebugData(&fields)
	l.log("error", message, fields)
}

func (l *loggerImpl) FatalFields(message string, fields map[string]interface{}) {
	l.appendDebugData(&fields)
	l.log("fatal", message, fields)
}

func (l *loggerImpl) PanicFields(message string, fields map[string]interface{}) {
	l.appendDebugData(&fields)
	l.log("panic", message, fields)
}

func appendCallerInfo(fields *map[string]interface{}) {
	_, file, line, ok := runtime.Caller(2)
	buf := make([]byte, 1<<16)
	stackSize := runtime.Stack(buf, true)

	if ok {
		if fields == nil {
			*fields = make(map[string]interface{})
		}

		(*fields)["file"] = file
		(*fields)["line"] = line
		(*fields)["stack"] = string(buf[0:stackSize])
	}
}
