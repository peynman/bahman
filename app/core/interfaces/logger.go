package interfaces

import "github.com/sirupsen/logrus"

type LoggingChannel interface {
	Config(string) bool
	GetLogger() *logrus.Logger
	GetChannelName() string
}