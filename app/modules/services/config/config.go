package config

import (
	"fmt"
	"github.com/peyman-abdi/bahman/app/interfaces/services"
	"github.com/peyman-abdi/conf"
)

type configImpl struct {
	config *conf.Config
}
func (c *configImpl) IsSet(key string) bool {
	return c.config.IsSet(key)
}
func (c *configImpl) Get(key string, def interface{}) interface{} {
	return c.config.Get(key, def)
}
func (c *configImpl) GetString(key string, def string) string {
	return c.config.GetString(key, def)
}
func (c *configImpl) GetInt(key string, def int) int {
	return c.config.GetInt(key, def)
}
func (c *configImpl) GetInt64(key string, def int64) int64 {
	return c.config.GetInt64(key, def)
}
func (c *configImpl) GetFloat(key string, def float64) float64 {
	return c.config.GetFloat(key, def)
}
func (c *configImpl) GetBoolean(key string, def bool) bool {
	return c.config.GetBoolean(key, def)
}
func (c *configImpl) GetStringArray(key string, def []string) []string {
	return c.config.GetStringArray(key, def)
}
func (c *configImpl) GetIntArray(key string, def []int) []int {
	return c.config.GetIntArray(key, def)
}
func (c *configImpl) GetFloatArray(key string, def []float64) []float64 {
	return c.config.GetFloatArray(key, def)
}
func (c *configImpl) GetMap(key string, def map[string]interface{}) map[string]interface{} {
	return c.config.GetMap(key, def)
}
func (c *configImpl) GetAsString(key string, def string) string {
	return c.config.GetAsString(key, def)
}


func New(app services.Application) services.Config {
	c, err := conf.New(app.ConfigPath(""), app.RootPath(""), []conf.EvaluatorFunction{
		NewStoragePathEvaluator(app),
		NewResourcesPathEvaluator(app),
		NewRootPathEvaluator(app),
		new(SystemParameterEvaluator),
		new(TimeEvaluator),
	})

	if err != nil && c == nil {
		panic(err)
	}
	if err != nil {
		fmt.Println(err)
	}

	confImpl := new(configImpl)
	confImpl.config = c

	return confImpl
}