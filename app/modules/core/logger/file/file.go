package main

import (
	"fmt"
	"github.com/peyman-abdi/avalanche/app/interfaces/core"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"time"
)

type AvalanchePluginFile struct {
}

func (_ *AvalanchePluginFile) Initialize(services core.Services) bool {
	fileLogger = new(FileLogChannel)
	fileLogger.services = services
	return true
}
func (_ *AvalanchePluginFile) Interface() interface{} {
	return fileLogger
}

var PluginInstance core.AvalanchePlugin = new(AvalanchePluginFile)

type FileLogChannel struct {
	services core.Services
	logger   *logrus.Logger
}

var _ core.LoggingChannel = (*FileLogChannel)(nil)
var fileLogger *FileLogChannel

func (logger *FileLogChannel) Config(conf_path string) bool {
	console := fileLogger
	console.logger = logrus.New()
	config := logger.services.Config()
	app := logger.services.App()

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
	console.logger.Out = file

	switch config.GetString(conf_path+".format", "text") {
	case "json":
		console.logger.Formatter = new(logrus.JSONFormatter)
	}

	setLoggerLevel(console.logger, config.GetString(conf_path+".level", "debug"))

	return true
}
func (logger *FileLogChannel) GetLogger() *logrus.Logger {
	return fileLogger.logger
}
func (logger *FileLogChannel) GetChannelName() string {
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
