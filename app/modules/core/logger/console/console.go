package main

import (
	"github.com/peyman-abdi/avalanche/app/interfaces/core"
	"github.com/sirupsen/logrus"
)

type AvalanchePluginConsole struct {
}

func (_ *AvalanchePluginConsole) Initialize(services core.Services) bool {
	consoleLogger = new(ConsoleLogChannel)
	consoleLogger.services = services
	return true
}
func (_ *AvalanchePluginConsole) Interface() interface{} {
	return consoleLogger
}

var PluginInstance core.AvalanchePlugin = new(AvalanchePluginConsole)

type ConsoleLogChannel struct {
	services core.Services
	logger   *logrus.Logger
}

var _ core.LoggingChannel = (*ConsoleLogChannel)(nil)
var consoleLogger *ConsoleLogChannel

func (logger *ConsoleLogChannel) Config(confPath string) bool {
	consoleLogger.logger = logrus.New()

	switch logger.services.Config().GetString(confPath+".format", "text") {
	case "json":
		consoleLogger.logger.Formatter = new(logrus.JSONFormatter)
	}

	setLoggerLevel(consoleLogger.logger, logger.services.Config().GetString(confPath+".level", "debug"))

	return true
}
func (logger *ConsoleLogChannel) GetLogger() *logrus.Logger {
	return consoleLogger.logger
}
func (logger *ConsoleLogChannel) GetChannelName() string {
	return "console"
}

func setLoggerLevel(logger *logrus.Logger, formatter string) {
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
