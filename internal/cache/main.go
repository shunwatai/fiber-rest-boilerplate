package cache

import (
	"fmt"
	"golang-api-starter/internal/config"
	logger "golang-api-starter/internal/helper/logger/zap_log"
	"strings"
	"sync"
)

var mu sync.RWMutex

type ICaching interface {
	// Get the caching client to its correspond struct
	SetClient() error

	// Get caching client info
	GetConnectionInfo() *ConnectionInfo

	// Get cache by key
	Get(key string, dst interface{}) bool

	// Set cache by key
	Set(key string, value interface{}) error

	// Delete cache by key
	DelByKey(key string) error

	// Delete all keys from cache
	FlushDb() error
}

var CacheService ICaching = nil

type ConnectionInfo struct {
	Driver   string
	Host     string
	Port     string
	User     *string
	Pass     *string
	Database *string
}

var cfg = config.Cfg

func NewCachingService() ICaching {
	if cfg.CacheConf == nil {
		logger.Errorf("error: DbConf is nil, maybe fail to load the config....")
	}
	logger.Debugf("engine: %+v", cfg.CacheConf.Driver)

	if cfg.CacheConf.Driver == "redis" {
		err := Rds.SetClient()
		if err != nil {
			logger.Fatalf("failed to initilise redis... err: %+v", err.Error())
		}

		return Rds
	}

	if cfg.CacheConf.Driver == "memcached" {
		err := Mc.SetClient()
		if err != nil {
			logger.Fatalf("failed to initilise memcached... err: %+v", err.Error())
		}

		return Mc
	}

	logger.Fatalf("failed to initilise caching service...")
	return nil
}

func GetCacheKey(key string, queryString map[string]interface{}) string {
	if len(key) == 0 {
		logger.Errorf("error: empty key....")
		return ""
	}

	if len(queryString) == 0 {
		return key
	}

	queries := []string{}
	for k, v := range queryString {
		if k == "exactMatch" || k == "columns" {
			continue
		}
		// logger.Debugf("key??>> %+v", k)
		switch v.(type) {
		case string:
			queries = append(queries, fmt.Sprintf("%s=%s", k, v))
		case []string:
			queries = append(queries, fmt.Sprintf("%s=%s", k, strings.Join(v.([]string), ",")))
		default:
		}
	}

	key = fmt.Sprintf("%s-%s", key, strings.Join(queries, "-"))
	return key
}

func EmptyCacheKeyMap(cachedKeys map[string]struct{}) {
	for k, _ := range cachedKeys {
		CacheService.DelByKey(k)
		delete(cachedKeys, k)
	}
}
