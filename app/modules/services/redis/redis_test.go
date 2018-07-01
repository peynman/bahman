package redis_test

import (
	"testing"
	"github.com/peyman-abdi/avest"
	"github.com/peyman-abdi/avalanche/app/modules/services/config"
	"github.com/peyman-abdi/avalanche/app/interfaces/services"
	application "github.com/peyman-abdi/avalanche/app/modules/services/app"
	redis2 "github.com/peyman-abdi/avalanche/app/modules/services/redis"
	"time"
)

var app = application.Initialize(0, "test")
var conf services.Config
var redis services.RedisClient
var configs = map[string]interface{}{
	"redis.hjson": map[string]interface{}{
		"client": map[string]interface{} {
			"host": "127.0.0.1",
			"port": "6379",
		},
	},
}
func init() {
	avest.CreateConfigFiles(app, configs)
	conf = config.Initialize(app)
	redis = redis2.Initialize(conf)
}

func TestRedisCommands(t *testing.T) {
	v, err := redis.Ping()
	if v != "PONG" || err != nil {
		t.Errorf("Ping failed: %s, %s", v, err)
	}

	expires := time.Now().Add(time.Second * 2)
	redis.Set("key", "something", &expires)

	if s, err := redis.Get("key"); s != "something" || err != nil {
		t.Errorf("Get failed: %s, %s", s, err)
	}

	time.Sleep(time.Second * 2)
	if s, err := redis.Get("key"); err != nil {
		t.Errorf("Get failed: %s, %s", s, err)
	}

	redis.Set("forever", "value", nil)
	if !redis.Exists("forever") {
		t.Errorf("key forever does not exist")
	}
	redis.Delete("forever")
	if redis.Exists("forever") {
		t.Errorf("key forever should be deleted")
	}

	if redis.Exists("not_found") {
		t.Errorf("key should not exist")
	}


}