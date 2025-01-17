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
		Driver: "redis",
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
	redisClient := r.Client

	err := redisClient.Get(ctx, key).Scan(dst)
	if err != nil {
		// logger.Errorf("failed to get cache, err: %+v", err.Error())
		return false
	}

	return true
}

func (r *Redis) Set(key string, value interface{}) error {
	redisClient := r.Client

	if err := redisClient.Set(ctx, key, value, 4*time.Hour /* 0 for no expire */).Err(); err != nil {
		return logger.Errorf("failed to set cache, err: %+v", err.Error())
	}

	return nil
}

func (r *Redis) DelByKey(key string) error {
	redisClient := r.Client
	redisClient.Del(ctx, key)

	return nil
}

// FlushDb for clear all keys for debug
func (r *Redis) FlushDb() error {
	redisClient := r.Client
	redisClient.FlushDB(ctx)

	return nil
}
