package main

import (
	"avalanche/app/core/interfaces"
	"github.com/sirupsen/logrus"
	"os"
	"fmt"
	"path/filepath"
	"time"
	"avalanche/app/core/config"
	"avalanche/app/core/app"
)


type AvalanchePluginFile struct {
}
func (_ *AvalanchePluginFile) Version() string {
	return "1.0.0"
}
func (_ *AvalanchePluginFile) VersionCode() int {
	return 1
}
func (_ *AvalanchePluginFile) AvalancheVersionCode() int {
	return app.VersionCode
}
func (_ *AvalanchePluginFile) Title() string {
	return "File Logger"
}
func (_ *AvalanchePluginFile) Description() string {
	return "File driver for Avalanche logger"
}
func (_ *AvalanchePluginFile) Initialize() bool {
	fileLogger = new(FileLogChannel)
	return true
}
func (_ *AvalanchePluginFile) Interface() interface{} {
	return fileLogger
}
var PluginInstance interfaces.AvalanchePlugin = new(AvalanchePluginFile)


type FileLogChannel struct {
	logger *logrus.Logger
}
var _ interfaces.LoggingChannel = (*FileLogChannel)(nil)
var fileLogger *FileLogChannel

func (logger *FileLogChannel) Config(conf_path string) bool {
	console := fileLogger
	console.logger = logrus.New()

	nowString := time.Now().Format("2006-01-02")
	logFilePath := config.GetString(conf_path + ".path", app.StoragePath("logs/" + nowString + ".log"))
	switch config.GetString(conf_path+".driver", "file") {
	case "daily":
		 logDir, logFilename := filepath.Split(logFilePath)
		 logFilenameExt := filepath.Ext(logFilename)
		 logFilename = logFilename[:len(logFilename) - len(logFilenameExt)]
		 logFilePath = filepath.Join(logDir, logFilename + "-" + nowString + ".log")
	default:

	}
	file, err := os.OpenFile(logFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0700)
	if err != nil {
		fmt.Println(err)
		return false
	}
	console.logger.Out = file

	switch config.GetString(conf_path + ".format", "text") {
	case "json":
		console.logger.Formatter = new(logrus.JSONFormatter)
	}

	setLoggerLevel(console.logger, config.GetString(conf_path + ".level", "debug"))

	return true
}
func (logger *FileLogChannel) GetLogger() *logrus.Logger {
	return fileLogger.logger
}
func (logger *FileLogChannel) GetChannelName() string {
	return "file"
}

func setLoggerLevel(logger *logrus.Logger, formatter string)  {
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
