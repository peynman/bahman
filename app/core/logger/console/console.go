package main

import (
	"github.com/sirupsen/logrus"
	"github.com/peyman-abdi/avalanche/app/interfaces"
	"github.com/peyman-abdi/avalanche/app/core/app"
)

type AvalanchePluginConsole struct {
}
func (_ *AvalanchePluginConsole) Version() string {
	return "1.0.0"
}
func (_ *AvalanchePluginConsole) VersionCode() int {
	return 1
}
func (_ *AvalanchePluginConsole) AvalancheVersionCode() int {
	return app.VersionCode
}
func (_ *AvalanchePluginConsole) Title() string {
	return "Console Logger"
}
func (_ *AvalanchePluginConsole) Description() string {
	return "Console driver for Avalanche logger"
}
func (_ *AvalanchePluginConsole) Initialize(services interfaces.Services) bool {
	consoleLogger = new(ConsoleLogChannel)
	consoleLogger.services = services
	return true
}
func (_ *AvalanchePluginConsole) Interface() interface{} {
	return consoleLogger
}
var PluginInstance interfaces.AvalanchePlugin = new(AvalanchePluginConsole)

type ConsoleLogChannel struct {
	services interfaces.Services
	logger *logrus.Logger
}
var _ interfaces.LoggingChannel = (*ConsoleLogChannel)(nil)
var consoleLogger *ConsoleLogChannel

func (logger *ConsoleLogChannel) Config(conf_path string) bool {
	consoleLogger.logger = logrus.New()

	switch logger.services.Config().GetString(conf_path + ".format", "text") {
	case "json":
		consoleLogger.logger.Formatter = new(logrus.JSONFormatter)
	}

	setLoggerLevel(consoleLogger.logger, logger.services.Config().GetString(conf_path + ".level", "debug"))

	return true
}
func (logger *ConsoleLogChannel) GetLogger() *logrus.Logger {
	return consoleLogger.logger
}
func (logger *ConsoleLogChannel) GetChannelName() string {
	return "console"
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