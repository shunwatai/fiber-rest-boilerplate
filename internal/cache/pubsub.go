package cache

import (
	logger "golang-api-starter/internal/helper/logger/zap_log"

	"github.com/redis/go-redis/v9"
)

type IPubSub interface {
	// Get the caching client to its correspond struct
	Pub(channelName, message string) error

	// Get caching client info
	Sub(channelName string) *redis.PubSub
}

var PubSubService IPubSub = nil

func NewPubSubService() IPubSub {
	if cfg.CacheConf == nil {
		logger.Errorf("error: DbConf is nil, maybe fail to load the config....")
	}
	logger.Debugf("engine: %+v", cfg.CacheConf.Driver)

	switch cfg.CacheConf.Driver {
	case "redis":
		err := Rds.SetClient()
		if err != nil {
			logger.Fatalf("failed to initilise redis... err: %+v", err.Error())
		}
		return Rds

	default:
		logger.Fatalf("failed to initilise pubsub service...")
		return nil
	}
}

func (r *Redis) Pub(channelName, message string) error {
	r.Client.Publish(ctx, channelName, message)
	logger.Debugf("published to %+v to boardcast", channelName)
	return nil
}

func (r *Redis) Sub(channelName string) *redis.PubSub {
	logger.Debugf("subscribed to redis, channel: %+v", channelName)
	return r.Client.Subscribe(ctx, channelName)
}

