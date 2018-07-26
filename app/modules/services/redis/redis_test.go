package redis_test

import (
	"testing"
	"github.com/peyman-abdi/avest"
	"github.com/peyman-abdi/bahman/app/modules/services/config"
	"github.com/peyman-abdi/bahman/app/interfaces/services"
	application "github.com/peyman-abdi/bahman/app/modules/services/app"
	redis2 "github.com/peyman-abdi/bahman/app/modules/services/redis"
	"time"
	"github.com/peyman-abdi/bahman/app/modules/services/logger"
	"gopkg.in/mgo.v2/bson"
)

var app = application.New(0, "test")
var conf services.Config
var log services.Logger
var redis services.RedisClient
var configs = map[string]interface{}{
	"redis.hjson": map[string]interface{}{
		"connections": map[string]interface{} {
			"local": map[string]interface{} {
				"host": "127.0.0.1",
				"port": "6379",
			},
		},
	},
}
func init() {
	avest.CreateConfigFiles(app, configs)
	conf = config.New(app)
	log = logger.New(conf)
	redis = redis2.New(conf, log)
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

	redis.Set("string_me", "string", nil)
	if str, err := redis.Get("string_me"); str != "string" {
		t.Error(err)
		t.Errorf("get as string failed")
	}

	testMap := map[string]interface{} {
		"one": 12,
		"two": "three",
	}
	testMapBin, err := bson.Marshal(testMap)
	if err != nil {
		t.Error(err)
	}
	if err = redis.Set("map", testMapBin, nil); err != nil {
		t.Error(err)
	}

	var scanBin string
	var scanMap map[string]interface{}
	if err = redis.Scan("map", &scanBin); err != nil {
		t.Error(err)
	}
	err = bson.Unmarshal([]byte(scanBin), &scanMap)
	if err != nil {
		t.Error(err)
	}
	if scanMap["one"] != 12 {
		t.Errorf("not the same object scan")
	}
}

func TestMoreRedisCommands(t *testing.T) {
	nestMap := map[string]interface{} {
		"nested": map[string]interface{} {
			"data": "some value",
		},
	}
	err := redis.Set("nested_map", nestMap, nil)
	if err != nil {
		t.Error(err)
	}
	var nestMapRef map[string]interface{}
	err = redis.Scan("nested_map", &nestMapRef)
	if err != nil {
		t.Error(err)
	}
	if nestMapRef["nested"].(map[string]interface{})["data"] != "some value" {
		t.Errorf("map is not the same %v, %v", nestMap, nestMapRef)
	}
}