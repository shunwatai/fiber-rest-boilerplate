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

func (mc *Memcached) Set(prefix,key string, value interface{}) error {
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

	if _, ok := cachedKeys[prefix]; !ok {
		cachedKeys[prefix] = []string{}
	}
	cachedKeys[prefix] = append(cachedKeys[prefix], key)

	return nil
}

func (mc *Memcached) DelByKey(key string) error {
	client := mc.Client
	return client.Delete(key)
}

func (mc *Memcached) DelByPrefix(prefix string) error {
	client := mc.Client
	if keys, ok := cachedKeys[prefix]; !ok {
		return logger.Errorf("Failed to get keys from cachedKeys[mc.KeyPrefix] by prefix: %s....", prefix)
	} else {
		for _, k := range keys {
			if err := client.Delete(k); err != nil {
				return logger.Errorf("Memcache failed to Delete key: %s....", k)
			}
		}
	}

	delete(cachedKeys, prefix) 
	return nil
}

// FlushDb for clear all keys for debug
func (mc *Memcached) FlushDb() error {
	client := mc.Client
	client.DeleteAll()

	return nil
}
