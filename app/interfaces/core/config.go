package core

type Config interface {
	IsSet(key string) bool
	Get(key string, def interface{}) interface{}
	GetString(key string, def string) string
	GetInt(key string, def int) int
	GetInt64(key string, def int64) int64
	GetFloat(key string, def float64) float64
	GetBoolean(key string, def bool) bool
	GetStringArray(key string, def []string) []string
	GetIntArray(key string, def []int) []int
	GetFloatArray(key string, def []float64) []float64
	GetMap(key string, def map[string]interface{}) map[string]interface{}
	GetAsString(key string, def string) string
}
