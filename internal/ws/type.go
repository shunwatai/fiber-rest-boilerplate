package ws

import (
	"encoding/json"
	"slices"
	"sync"

	"golang-api-starter/internal/cache"
	"golang-api-starter/internal/config"
	logger "golang-api-starter/internal/helper/logger/zap_log"
	"golang-api-starter/internal/modules/groupUser"
	"golang-api-starter/internal/modules/user"
)

var cfg = config.Cfg

type IOnlineUserList interface {
	Get(key string, dst *groupUser.User) bool
	Set(key string, value *groupUser.User)
	Del(key string)
	GetList() groupUser.Users
}

// onlineUserListMap for local without 3rd party cache service like redis
type onlineUserListMap struct {
	list sync.Map
}

func (oulm *onlineUserListMap) Get(key string, dst *groupUser.User) bool {
	val, ok := oulm.list.Load(key)
	dst = val.(*groupUser.User)
	return ok
}

func (oulm *onlineUserListMap) Set(key string, value *groupUser.User) {
	oulm.list.Store(key, value)
}

func (oulm *onlineUserListMap) Del(key string) {
	oulm.list.Delete(key)
}

func (oulm *onlineUserListMap) GetList() groupUser.Users {
	userList := groupUser.Users{}
	oulm.list.Range(func(key, value interface{}) bool {
		userList = append(userList, value.(*groupUser.User))
		return true // continue iteration
	})
	return userList
}

// onlineUserListRds for storing the userList in cache service like redis
type onlineUserListRds struct {
	keys         cachedKeys
	list         cache.ICaching
	pubsub       cache.IPubSub
	listCacheKey string
}

type cachedKeys []string

func (s cachedKeys) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}

// make sure the Student interface here accepts a pointer
func (s *cachedKeys) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, s)
}

func (oulm *onlineUserListRds) Get(key string, dst *groupUser.User) bool {
	cacheVal := &user.CacheValue{Users: []*groupUser.User{dst}}
	ok := oulm.list.Get(key, cacheVal)
	return ok
}

func (oulm *onlineUserListRds) Set(key string, value *groupUser.User) {
	defer func() {
		if err := oulm.list.Set(oulm.listCacheKey, &oulm.keys); err != nil {
			logger.Errorf("failed to set keys: %+v to cache..., err: %+v", oulm.listCacheKey, err.Error())
		}
	}()
	if ok := oulm.list.Get(oulm.listCacheKey, &oulm.keys); !ok {
		logger.Errorf("failed to get key: %+v from cache...", oulm.listCacheKey)
	}
	oulm.keys = append(oulm.keys, key)
	if err := oulm.list.Set(key, &user.CacheValue{Users: []*groupUser.User{value}}); err != nil {
		logger.Errorf("failed to set key: %+v to cache...", key)
	}
}

func (oulm *onlineUserListRds) Del(key string) {
	defer func() {
		if err := oulm.list.Set(oulm.listCacheKey, &oulm.keys); err != nil {
			logger.Errorf("failed to set keys: %+v to cache..., err: %+v", oulm.listCacheKey, err.Error())
		}
	}()
	if ok := oulm.list.Get(oulm.listCacheKey, &oulm.keys); !ok {
		logger.Errorf("failed to get key: %+v from cache...", oulm.listCacheKey)
	}
	keyIdx := -1
	for i, k := range oulm.keys {
		if k == key {
			keyIdx = i
			break
		}
	}
	oulm.keys = slices.Delete(oulm.keys, keyIdx, keyIdx+1)

	if err := oulm.list.DelByKey(key); err != nil {
		logger.Errorf("failed to del key: %+v in cache...", key)
	}
}

func (oulm *onlineUserListRds) GetList() groupUser.Users {
	if ok := oulm.list.Get(oulm.listCacheKey, &oulm.keys); !ok {
		logger.Errorf("failed to get key: %+v from cache...", oulm.listCacheKey)
	}
	users := []*groupUser.User{}

	logger.Debugf("oulm.keys: %+v", oulm.keys)
	for _, k := range oulm.keys {
		u := groupUser.User{}
		oulm.Get(k, &u)
		users = append(users, &u)
	}
	return users
}

func NewOnlineUserList() IOnlineUserList {
	if cfg.CacheConf.Enabled {
		logger.Debugf("use redis for userlist")
		return &onlineUserListRds{
			list:         cache.CacheService,
			pubsub:       cache.PubSubService,
			keys:         []string{},
			listCacheKey: "user-list-keys",
		} // return redis
	}

	logger.Debugf("use local map for userlist")
	// without redis
	return &onlineUserListMap{
		list: sync.Map{},
	}
}
