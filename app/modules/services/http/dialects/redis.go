package dialects

import (
	"github.com/peyman-abdi/bahman/app/interfaces/services"
	"github.com/fasthttp-contrib/sessions"
)

type RedisSessionDriver struct {
	redis services.RedisClient
	logger services.Logger
}

func (r *RedisSessionDriver) Load(key string) map[string]interface{} {
	var m map[string]interface{}
	err := r.redis.Scan(key, m)
	if err != nil {
		r.logger.ErrorFields("RedisSessionDriver failed loading key", map[string]interface{} {
			"err": err,
			"key": key,
		})
		return nil
	}
	return m
}

func (r *RedisSessionDriver) Update(key string, obj map[string]interface{}) {
	r.redis.Set(key, obj, nil)
}

func NewRedisDriver(redis services.RedisClient, logger services.Logger) sessions.Database {
	return &RedisSessionDriver{
		redis: redis,
		logger: logger,
	}
}
