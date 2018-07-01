package redis

import (
	"github.com/peyman-abdi/avalanche/app/interfaces/services"
	"github.com/go-redis/redis"
	"time"
)

type redisImpl struct {
	cluster *redis.ClusterClient
	client *redis.Client
}

var instance *redisImpl
var _ services.RedisClient = (*redisImpl)(nil)

func Initialize(config services.Config) services.RedisClient {
	instance = new(redisImpl)

	clusters := config.GetStringArray("redis.clusters", []string{})
	if clusters != nil && len(clusters) > 0 {
		instance.cluster = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs: clusters,
			Password: config.GetAsString("redis.client.password", ""),
		})
	} else {
		instance.client = redis.NewClient(&redis.Options{
			Addr: config.GetString("redis.client.host", "") + ":" +  config.GetAsString("redis.client.port", "6379"),
			Password: config.GetAsString("redis.client.password", ""),
			DB: config.GetInt("redis.client.db", 0),
		})
	}

	return instance
}

func Close() {
	if instance.client != nil {
		instance.client.Close()
	}
	if instance.cluster != nil {
		instance.cluster.Close()
	}
}

func (r *redisImpl) Ping() (string, error) {
	return r.client.Ping().Result()
}

func (r *redisImpl) Get(key string) (interface{}, error)  {
	if r.Exists(key) {
		return r.client.Get(key).Result()
	}

	return nil, nil
}

func (r *redisImpl) GetAsString(key string) (string, error)  {
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
	cmd := r.client.Get(key)
	if cmd.Err() != nil {
		return cmd.Err()
	}

	if err := cmd.Scan(target); err != nil {
		return err
	}

	return nil
}
func (r *redisImpl) Set(key string, val interface{}, expiration *time.Time) error {
	var d int64 = 0
	if expiration != nil {
		d = expiration.UnixNano() - time.Now().UnixNano()
	}
	return r.client.Set(key, val, time.Duration(d)).Err()
}