package services

import (
	"time"
)

type Cache interface {
	Forever(key string, val interface{})
	Set(key string, val interface{}, expires *time.Time)
	Get(key string) interface{}
	Exists(key string) bool
	Add(key string, val interface{}, expires *time.Time) bool
	Delete(key string)
}

