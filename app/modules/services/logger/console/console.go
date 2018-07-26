package console

import (
	"github.com/peyman-abdi/bahman/app/interfaces/services"
	"github.com/sirupsen/logrus"
)

type LogChannel struct {
	services services.Services
	logger   *logrus.Logger
}

var _ services.LoggingChannel = (*LogChannel)(nil)

func (c *LogChannel) Config(confPath string) bool {
	c.logger = logrus.New()

	switch c.services.Config().GetString(confPath+".format", "text") {
	case "json":
		c.logger.Formatter = new(logrus.JSONFormatter)
	}

	setLoggerLevel(c.logger, c.services.Config().GetString(confPath+".level", "debug"))

	return true
}
func (c *LogChannel) GetLogger() interface{} {
	return c.logger
}
func (c *LogChannel) GetChannelName() string {
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

func New(instance services.Services) *LogChannel {
	return &LogChannel{
		services: instance,
	}
}