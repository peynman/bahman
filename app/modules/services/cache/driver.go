package cache

import "time"

type DriverDialect interface {
	Get(string) (interface{}, error)
	Set(key string, val interface{}, expires *time.Time) error
	Exists(key string) bool
	Delete(key string) error
}
