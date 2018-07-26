package services

import "time"

type RedisClient interface {
	Connection(name string) RedisClient

	Ping() (string, error)
	Set(key string, value interface{}, expiration *time.Time) error
	Get(key string) (string, error)
	Scan(key string, target interface{}) error
	Exists(key string) bool
	Delete(key string) error
	GetAsMap(key string) (map[string]interface{},error)
	GetAsArray(key string) ([]interface{},error)
}
