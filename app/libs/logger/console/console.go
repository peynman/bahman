package main

import (
	"github.com/sirupsen/logrus"
	"avalanche/app/libs/interfaces"
	"AvConfig/lib"
)


type ConsoleLogChannel struct {
	logger *logrus.Logger
}
var _ interfaces.LoggingChannel = (*ConsoleLogChannel)(nil)

var consoleLogger interfaces.LoggingChannel = & ConsoleLogChannel {
	logrus.New(),
}

func (logger *ConsoleLogChannel) Config(conf_path string) bool {

	console := consoleLogger.(*ConsoleLogChannel)

	setLoggerLevel(console.logger, config.GetString(conf_path + ".level", "debug"))

	return true
}
func (logger *ConsoleLogChannel) GetLogger() logrus.Logger {
	return *(consoleLogger.(*ConsoleLogChannel)).logger
}

func setLoggerLevel(logger *logrus.Logger, formatter string)  {
	switch formatter {
	case "debug":
		logger.SetLevel(logrus.DebugLevel)
	case "info":
		logger.SetLevel(logrus.InfoLevel)
	case "warn":
		logger.SetLevel(logrus.WarnLevel)
	case "error":
		logger.SetLevel(logrus.ErrorLevel)
	case "fatal":
		logger.SetLevel(logrus.FatalLevel)
	case "panic":
		logger.SetLevel(logrus.PanicLevel)
	}
}

func GetInstance() interfaces.LoggingChannel {
	return consoleLogger
}