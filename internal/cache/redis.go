package cache

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"

	logger "golang-api-starter/internal/helper/logger/zap_log"
)

type Redis struct {
	Client *redis.Client
}

var Rds = &Redis{}

// GetConnectionInfo get cache's var by config
func (r *Redis) GetConnectionInfo() *ConnectionInfo {
	cfg.LoadEnvVariables()
	return &ConnectionInfo{
		Driver: cfg.CacheConf.Driver,
		Host:   cfg.CacheConf.RedisConf.Host,
		Port:   cfg.CacheConf.RedisConf.Port,
		User:   nil,
		Pass:   nil,
	}
}

// GetConnOption returns the redis options
func (r *Redis) GetConnOption() *redis.Options {
	connInfo := Rds.GetConnectionInfo()
	rdhost := connInfo.Host
	rdport := connInfo.Port

	rdOptions := &redis.Options{
		Addr:     fmt.Sprintf("%s:%s", rdhost, rdport),
		Password: "", // no password set
		DB:       0,  // use default DB
	}

	return rdOptions
}

// SetClient initiate the redis client by GetConnOption
func (r *Redis) SetClient() error {
	var client = redis.NewClient(r.GetConnOption())
	var ctx = context.Background()

	if err := client.Ping(ctx).Err(); err != nil {
		return err
	}

	r.Client = client

	return nil
}

var ctx = context.Background()

func (r *Redis) Get(key string, dst interface{}) bool {
	client := r.Client

	err := client.Get(ctx, key).Scan(dst)
	if err != nil {
		// logger.Errorf("failed to get cache, err: %+v", err.Error())
		return false
	}

	// logger.Debugf(">>> cachedKeys: %+v", cachedKeys)
	return true
}

func (r *Redis) Set(prefix, key string, value interface{}) error {
	client := r.Client

	if err := client.Set(ctx, key, value, 4*time.Hour /* 0 for no expire */).Err(); err != nil {
		return logger.Errorf("failed to set cache, err: %+v", err.Error())
	}

	if _, ok := cachedKeys[prefix]; !ok {
		cachedKeys[prefix] = []string{}
	}

	cachedKeys[prefix] = append(cachedKeys[prefix], key)
	// logger.Debugf("RDS set key?? %+v, %+v", prefix, cachedKeys)

	return nil
}

func (r *Redis) GetKeysByPrefix(prefix string) ([]string, error) {
	keys := []string{}
	pattern := fmt.Sprintf("%s:*", prefix)
	iter := r.Client.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}
	if err := iter.Err(); err != nil {
		return keys, err
	}

	return keys, nil
}

func (r *Redis) DelByKey(key string) error {
	client := r.Client
	client.Del(ctx, key)

	return nil
}

func (r *Redis) DelByPrefix(prefix string) error {
	client := r.Client

	if keys, err := r.GetKeysByPrefix(prefix); err != nil {
		// logger.Errorf("GetKeysByPrefix, err: %+v", err.Error())
	} else {

		for _, k := range keys {
			// logger.Debugf("Deleting key: %+v", k)
			client.Del(ctx, k)
		}
	}

	delete(cachedKeys, prefix) 
	return nil
}

// FlushDb for clear all keys for debug
func (r *Redis) FlushDb() error {
	client := r.Client
	client.FlushDB(ctx)

	return nil
}
