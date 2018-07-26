package file

import (
	"fmt"
	"github.com/peyman-abdi/bahman/app/interfaces/services"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"time"
)

type LogChannel struct {
	services services.Services
	logger   *logrus.Logger
}

var _ services.LoggingChannel = (*LogChannel)(nil)

func (c *LogChannel) Config(conf_path string) bool {
	c.logger = logrus.New()
	config := c.services.Config()
	app := c.services.App()

	nowString := time.Now().Format("2006-01-02")
	logFilePath := config.GetString(conf_path+".path", app.StoragePath("logs/"+nowString+".log"))
	switch config.GetString(conf_path+".driver", "file") {
	case "daily":
		logDir, logFilename := filepath.Split(logFilePath)
		logFilenameExt := filepath.Ext(logFilename)
		logFilename = logFilename[:len(logFilename)-len(logFilenameExt)]
		logFilePath = filepath.Join(logDir, logFilename+"-"+nowString+".log")
	default:

	}
	file, err := os.OpenFile(logFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0700)
	if err != nil {
		fmt.Println(err)
		return false
	}
	c.logger.Out = file

	switch config.GetString(conf_path+".format", "text") {
	case "json":
		c.logger.Formatter = new(logrus.JSONFormatter)
	}

	setLoggerLevel(c.logger, config.GetString(conf_path+".level", "debug"))

	return true
}
func (c *LogChannel) GetLogger() interface{} {
	return c.logger
}
func (c *LogChannel) GetChannelName() string {
	return "file"
}

func setLoggerLevel(logger *logrus.Logger, formatter string) {
	switch formatter {
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	case "warn":
		logrus.SetLevel(logrus.WarnLevel)
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
	case "fatal":
		logrus.SetLevel(logrus.FatalLevel)
	case "panic":
		logrus.SetLevel(logrus.PanicLevel)
	}
}

func New(instance services.Services) *LogChannel {
	return &LogChannel {
		services: instance,
	}
}