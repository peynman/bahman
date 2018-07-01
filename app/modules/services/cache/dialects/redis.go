package dialects

import (
	"time"
	"github.com/peyman-abdi/avalanche/app/interfaces/services"
)

type RedisCacheDialect struct {
	Client services.RedisClient
}

func (r *RedisCacheDialect) Get(key string) (interface{}, error) {
	return r.Client.Get(key)
}

func (r *RedisCacheDialect) Exists(key string) bool {
	return r.Client.Exists(key)
}

func (r *RedisCacheDialect) Set(key string, val interface{}, expires *time.Time) error {
	return r.Client.Set(key, val, expires)
}

func (r *RedisCacheDialect) Delete(key string) error {
	return r.Client.Delete(key)
}
