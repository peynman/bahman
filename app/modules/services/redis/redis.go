package redis

import (
	"github.com/peyman-abdi/bahman/app/interfaces/services"
	"github.com/go-redis/redis"
	"time"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"errors"
	"reflect"
)

type redisImpl struct {
	client redis.Cmdable
	connections map[string]redis.Cmdable
	logger services.Logger
}

func (r *redisImpl) GetAsMap(key string) (map[string]interface{}, error) {
	var scanBin string
	var scanMap map[string]interface{}
	if err := r.Scan(key, &scanBin); err != nil {
		return nil, err
	}
	err := bson.Unmarshal([]byte(scanBin), &scanMap)
	if err != nil {
		return nil, err
	}
	return scanMap, nil
}

func (r *redisImpl) GetAsArray(key string) ([]interface{}, error) {
	var scanBin string
	var scanArr []interface{}
	if err := r.Scan(key, &scanBin); err != nil {
		return nil, err
	}
	err := bson.Unmarshal([]byte(scanBin), &scanArr)
	if err != nil {
		return nil, err
	}
	return scanArr, nil
}

var _ services.RedisClient = (*redisImpl)(nil)
func (r *redisImpl) Connection(name string) services.RedisClient {
	if r.connections[name] == nil  {
		r.logger.ErrorFields("Redis Connection not found", map[string]interface{} {
			"error": fmt.Sprintf("Connection with name %s not found", name),
			"name": name,
			"connections": r.connections,
		})
		return nil
	}

	if r.connections[name] != nil {
		return &redisImpl{
			client: r.connections[name],
			connections: r.connections,
		}
	}

	return nil
}

func (r *redisImpl) Ping() (string, error) {
	return r.client.Ping().Result()
}

func (r *redisImpl) Get(key string) (string, error)  {
	if r.Exists(key) {
		return r.client.Get(key).Result()
	}

	return "", nil
}

func (r *redisImpl) Exists(key string) bool {
	e, err := r.client.Exists(key).Result()
	if err != nil {
		return false
	}
	return e != 0
}

func (r *redisImpl) Delete(key string) error {
	return r.client.Del(key).Err()
}

func (r *redisImpl) Scan(key string, target interface{}) error  {
	if !r.Exists(key) {
		return nil
	}

	switch target.(type) {
	case *map[string]interface{}, *[]interface{}:
		bin, err := r.client.Get(key).Result()
		if err != nil {
			return err
		}
		return bson.Unmarshal([]byte(bin), target)
	default:
		return r.client.Get(key).Scan(target)
	}
	return nil
}

func (r *redisImpl) Set(key string, val interface{}, expiration *time.Time) error {
	var d int64 = 0
	if expiration != nil {
		d = expiration.UnixNano() - time.Now().UnixNano()
	}
	switch val.(type) {
	case map[string]interface{}, []interface{}:
		bin, err := bson.Marshal(val)
		if err != nil {
			return err
		}
		return r.client.Set(key, bin, time.Duration(d)).Err()
	default:
		return r.client.Set(key, val, time.Duration(d)).Err()
	}

	return errors.New(fmt.Sprintf("unknown type of val interface: %s", reflect.TypeOf(val)))
}

func New(config services.Config, logger services.Logger) services.RedisClient {
	instance := new(redisImpl)
	instance.connections = make(map[string]redis.Cmdable)
	instance.logger = logger

	connectionConfigs := config.GetMap("redis.connections", map[string]interface{}{
		"local": map[string]interface{} {
			"host": "127.0.0.1",
			"port": 6379,
		},
	})
	for conn, params := range connectionConfigs {
		conconf := params.(map[string]interface{})
		if conconf["cluster"] == true {
			c := redis.NewClusterClient(&redis.ClusterOptions{
				Addrs: config.GetStringArray("redis.connections." + conn + ".clusters", []string{}),
				Password: config.GetAsString("redis.connections." + conn + ".password", ""),
			})
			instance.connections[conn] = c
		} else {
			c := redis.NewClient(&redis.Options{
				Addr: config.GetString("redis.connections." + conn + ".host", "") + ":" +  config.GetAsString("redis.connections." + conn + ".port", "6379"),
				Password: config.GetAsString("redis.connections." + conn + ".password", ""),
				DB: config.GetInt("redis.connections." + conn + ".db", 0),
			})
			instance.connections[conn] = c
		}
	}

	def := config.GetString("redis.default", "local")
	if instance.connections[def] == nil {
		panic(fmt.Sprintf("Default connection for Redis client not found: %s, %v", def, connectionConfigs))
	} else {
		instance.client = instance.connections[def]
		instance.connections["default"] = instance.client
	}

	return instance
}
