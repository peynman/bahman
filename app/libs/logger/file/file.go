package main

import (
	"avalanche/app/libs/interfaces"
	"github.com/sirupsen/logrus"
	"AvConfig/lib"
)

type FileLogChannel struct {
	logger *logrus.Logger
}
var _ interfaces.LoggingChannel = (*FileLogChannel)(nil)

var fileLogger interfaces.LoggingChannel = &FileLogChannel{
	logrus.New(),
}

func (logger *FileLogChannel) Config(conf_path string) bool {
	console := fileLogger.(*FileLogChannel)

	//file, err := os.Create()
	//if err != nil {
	//	fmt.Println(err)
	//	return falseåå
	//}
	//
	//console.logger.Out = file

	setLoggerLevel(console.logger, config.GetString(conf_path + ".level", "debug"))

	return true
}
func (logger *FileLogChannel) GetLogger() logrus.Logger {
	return *(fileLogger.(*FileLogChannel)).logger
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

func GetInstance() interfaces.LoggingChannel {
	return fileLogger
}
