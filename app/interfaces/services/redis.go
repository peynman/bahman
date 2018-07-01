package services

import "time"

type RedisClient interface {
	Ping() (string, error)
	Set(key string, value interface{}, expiration *time.Time) error
	Get(key string) (interface{}, error)
	GetAsString(key string) (string, error)
	Scan(key string, target interface{}) error
	Exists(key string) bool
	Delete(key string) error
}
