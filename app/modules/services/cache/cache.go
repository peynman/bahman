package cache

import (
	"github.com/peyman-abdi/avalanche/app/interfaces/services"
	"github.com/peyman-abdi/avalanche/app/modules/services/cache/dialects"
)

func Initialize(app services.Application, conf services.Config, log services.Logger, client services.RedisClient) services.Cache {
	c := new(cacheImpl)
	c.logger = log

	switch conf.GetAsString("cache.driver", "redis") {
	case "redis":
		c.driver = &dialects.RedisCacheDialect{ Client: client }
	}

	return c
}

