package config_test

import (
	"github.com/peyman-abdi/avalanche/app/interfaces/core"
	application "github.com/peyman-abdi/avalanche/app/modules/core/app"
	"github.com/peyman-abdi/avalanche/app/modules/core/config"
	"github.com/peyman-abdi/testil"
	"testing"
	"time"
)

var app = application.Initialize(0, "test")
var conf core.Config
var configs = map[string]interface{}{
	"simple.hjson": map[string]interface{}{
		"intTest":   1,
		"strTest":   "something",
		"floatTest": 1.3,
		"boolTest":  true,
	},
	"other.arrays.hjson": map[string]interface{}{
		"ints":   []int{1, 2, 3, 4, 5},
		"floats": []float64{1.1, 1.2, 1.3, 1.4},
	},
	"folder.evals.hjson": map[string]interface{}{
		"res":     "resources('test')",
		"storage": "storage('test')",
		"root":    "root()",
		"pl":      "system(\"platform\")",
		"vr":      "system(\"variant\")",
		"tm":      "time(\"now\")",
	},
	"maps.hjson": map[string]interface{}{
		"object": map[string]interface{}{
			"strTest": "something",
		},
	},
}

func init() {
	testil.CreateConfigFiles(app, configs)
	conf = config.Initialize(app)
}

func TestSimple(t *testing.T) {
	testString("simple.strTest", "something", t)
	testInt("simple.intTest", 1, t)
	testFloat("simple.floatTest", 1.3, t)
	testMapString("maps.object", "strTest", "something", t)
	testBoolean("simple.boolTest", true, t)
	testAsString("simple.boolTest", "true", t)
	testAsString("other.arrays.ints[1]", "2", t)
	testAsString("other.arrays.floats[1]", "1.2", t)
}
func TestEvaluations(t *testing.T) {
	testString("folder.evals.res", app.ResourcesPath("test"), t)
	testString("folder.evals.storage", app.StoragePath("test"), t)
	testString("folder.evals.root", app.RootPath(""), t)
	testString("folder.evals.pl", app.Platform(), t)
	testString("folder.evals.vr", app.Variant(), t)

	if v, ok := conf.Get("folder.evals.tm", 0).(time.Time); !ok {
		t.Errorf("Expecting time form %s, Got: %v", "folders.evals.tm", v)
	}
}

func testString(key string, value interface{}, t *testing.T) {
	if v := conf.GetString(key, "not found"); v != value {
		t.Error("Wrong value return from config key: " + key + "\nExpecting: " + value.(string) + "\nGot: " + v)
	}
}
func testInt(key string, value interface{}, t *testing.T) {
	if v := conf.GetInt(key, 0); v != value.(int) {
		t.Errorf("Wrong value return from config key: %s\nExpecting: %v\nGot: %v", key, value, v)
	}
}
func testFloat(key string, value interface{}, t *testing.T) {
	if v := conf.GetFloat(key, 0); v != value.(float64) {
		t.Errorf("Wrong value return from config key: %s\nExpecting: %v\nGot: %v", key, value, v)
	}
}
func testBoolean(key string, value interface{}, t *testing.T) {
	if v := conf.GetBoolean(key, false); v != value.(bool) {
		t.Errorf("Wrong value return from config key: %s\nExpecting: %v\nGot: %v", key, value, v)
	}
}
func testAsString(key string, value interface{}, t *testing.T) {
	if v := conf.GetAsString(key, "not found"); v != value.(string) {
		t.Errorf("Wrong value return from config key: %s\nExpecting: %v\nGot: %v", key, value, v)
	}
}
func testMapString(key string, innter string, value interface{}, t *testing.T) {
	if m := conf.GetMap(key, map[string]interface{}{}); m[innter] != value.(string) {
		t.Errorf("Wrong value return from config key: %s\nExpecting: %v\nGot: %v", key, value, m)
	}
}
