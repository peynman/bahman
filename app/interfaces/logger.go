package interfaces

import "github.com/sirupsen/logrus"

type LoggingChannel interface {
	Config(string) bool
	GetLogger() *logrus.Logger
	GetChannelName() string
}

type Logger interface {
	LoadConsole()
	LoadChannels(services Services)
    Debug(message string)
    Info(message string)
    Warn(message string)
    Error(message string)
    Fatal(message string)
    Panic(message string)
    DebugFields(message string, fields map[string]interface{})
    InfoFields(message string, fields map[string]interface{})
    WarnFields(message string, fields map[string]interface{})
    ErrorFields(message string, fields map[string]interface{})
    FatalFields(message string, fields map[string]interface{})
    PanicFields(message string, fields map[string]interface{})
}