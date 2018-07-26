package cache

import (
	"time"
	"github.com/peyman-abdi/bahman/app/interfaces/services"
)

type cacheImpl struct {
	driver DriverDialect
	logger services.Logger
}

func (c *cacheImpl) Forever(key string, val interface{}) {
	err := c.driver.Set(key, val, nil)
	if err != nil {

	}
}

var _ services.Cache = (*cacheImpl)(nil)

func (c *cacheImpl) Set(key string, val interface{}, expires *time.Time) {
	err := c.driver.Set(key, val, expires)
	if err != nil {
		c.logger.ErrorFields("Cache Set Error", map[string]interface {} {
			"error": err,
			"key": key,
			"value": val,
			"expires": expires,
		})
	}
}

func (c *cacheImpl) Get(key string) interface{} {
	 val, err := c.driver.Get(key)
	 if err != nil {
	 	c.logger.ErrorFields("Cache Get Error", map[string]interface {} {
	 		"error": err,
	 		"key": key,
		})
	 	return nil
	 }
	 return val
}

func (c *cacheImpl) Exists(key string) bool {
	return c.driver.Exists(key)
}

func (c *cacheImpl) Add(key string, val interface{}, expires *time.Time) bool {
	added := true
	if c.Exists(key) {
		added = false
	}

	c.Set(key, val, expires)
	return added
}

func (c *cacheImpl) Delete(key string) {
	c.driver.Delete(key)
}

