package config

import (
	"peyman/config"
	"avalanche/app/core/app"
)

var config *nemo.Config

func Initialize() {
	config = nemo.New(app.ConfigPath(""), app.RootPath(""), []nemo.EvaluatorFunction {
		new (StoragePathEvaluator),
		new (ResourcesPathEvaluator),
		new (RootPathEvaluator),
		new (SystemParameterEvaluator),
		new (TimeEvaluator),
	})
}

func  IsSet(key string) bool {
	return config.IsSet(key)
}
func  Get(key string, def interface{}) interface{} {
	return config.Get(key, def)
}
func  GetString(key string, def string) string {
	return config.GetString(key, def)
}
func  GetInt(key string, def int) int {
	return config.GetInt(key, def)
}
func  GetInt64(key string, def int64) int64 {
	return config.GetInt64(key, def)
}
func  GetFloat(key string, def float64) float64 {
	return config.GetFloat(key, def)
}
func  GetBoolean(key string, def bool) bool {
	return config.GetBoolean(key, def)
}
func  GetStringArray(key string, def []string) []string {
	return config.GetStringArray(key, def)
}
func  GetIntArray(key string, def []int) []int {
	return config.GetIntArray(key, def)
}
func  GetFloatArray(key string, def []float64) []float64 {
	return config.GetFloatArray(key, def)
}
func  GetMap(key string, def map[string]interface{}) map[string]interface{} {
	return config.GetMap(key, def)
}
func  GetAsString(key string, def string) string {
	return config.GetAsString(key, def)
}

