package cache

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	logger "golang-api-starter/internal/helper/logger/zap_log"
)

type Memcached struct {
	Client *memcache.Client
}

var Mc = &Memcached{}

// GetConnectionInfo get cache's var by config
func (mc *Memcached) GetConnectionInfo() *ConnectionInfo {
	cfg.LoadEnvVariables()
	return &ConnectionInfo{
		Driver: cfg.CacheConf.Driver,
		Host:   cfg.CacheConf.MemcachedConf.Host,
		Port:   cfg.CacheConf.MemcachedConf.Port,
		User:   nil,
		Pass:   nil,
	}
}

// SetClient initiate the memcached client by GetConnOption
func (mc *Memcached) SetClient() error {
	connInfo := Mc.GetConnectionInfo()
	host := connInfo.Host
	port := connInfo.Port
	client := memcache.New(fmt.Sprintf("%s:%s", host, port))

	if err := client.Ping(); err != nil {
		return err
	}

	mc.Client = client

	return nil
}

func (mc *Memcached) Get(key string, dst interface{}) bool {
	client := mc.Client

	mcItem, err := client.Get(key)
	if err != nil {
		// logger.Errorf("failed to get cache, err: %+v", err.Error())
		return false
	}

	if err := json.Unmarshal(mcItem.Value, &dst); err != nil {
		logger.Errorf("failed to Unmarshal cache, err: %+v", err.Error())
		return false
	}

	return true
}

func (mc *Memcached) Set(key string, value interface{}) error {
	client := mc.Client

	b, err := json.Marshal(&value)
	if err != nil {
		return logger.Errorf("failed to Marshal cache, err: %+v", err.Error())
	}

	if err := client.Set(&memcache.Item{
		Key:        key,
		Value:      b,
		Expiration: int32((4 * time.Hour).Seconds()),
	}); err != nil {
		return logger.Errorf("failed to set cache, err: %+v", err.Error())
	}

	return nil
}

func (mc *Memcached) DelByKey(key string) error {
	client := mc.Client
	return client.Delete(key)
}

// FlushDb for clear all keys for debug
func (mc *Memcached) FlushDb() error {
	client := mc.Client
	client.DeleteAll()

	return nil
}
