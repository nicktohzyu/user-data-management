package cache

import (
	"encoding/json"
	"github.com/gomodule/redigo/redis"
	"time"
	"user-data-management/commons"
	"user-data-management/commons/logger"
	"user-data-management/commons/user"
)

type RedisCache struct {
	dbNum    int
	duration time.Duration
	pool     *redis.Pool
}

const (
	ADDRESS          = "localhost:6379"
	EXPIRY_TIME_SECS = 60
)

func NewRedisCache() *RedisCache {
	cache := RedisCache{}
	cache.initPool()
	cache.ping()
	return &cache
}

func (cache *RedisCache) initPool() {
	cache.pool = &redis.Pool{
		MaxIdle:   80,
		MaxActive: 12000,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", ADDRESS)
			if err != nil {
				logger.Error("ERROR: fail init redis pool: ", err.Error())
				panic("Redis cache failed ping")
			}
			return conn, err
		},
	}
}

func (cache *RedisCache) ping() {
	conn := cache.pool.Get()
	defer conn.Close()
	_, err := redis.String(conn.Do("PING"))
	if err != nil {
		logger.Error("ERROR: fail ping redis conn: ", err.Error())
		panic("Redis cache failed ping")
	}
}

func (cache *RedisCache) Store(user1 user.User) error {
	startTime := time.Now()
	conn := cache.pool.Get()
	defer conn.Close()

	userStr, err := json.Marshal(user1)
	if err != nil {
		logger.Error("Marshal outgoing packet error")
		return err
	}

	_, err = conn.Do("SET", user1.Username, userStr)
	if err != nil {
		logger.Error("Cache failed to store user", err.Error())
		return err
	}

	_, err = conn.Do("EXPIRE", user1.Username, EXPIRY_TIME_SECS)
	if err != nil {
		logger.Error("Cache failed to set expiry time", err.Error())
		return err
	}

	endTime := time.Now()
	commons.Latency.WithLabelValues(
		commons.CacheComponent, commons.CacheStoreLabel).
		Observe(float64(endTime.Sub(startTime).Milliseconds()))
	return nil
}

func (cache *RedisCache) Get(username string) (*user.User, error) {
	startTime := time.Now()
	conn := cache.pool.Get()
	defer conn.Close()

	userStr, err := redis.String(conn.Do("GET", username))
	if err == redis.ErrNil {
		logger.Info("User not found in cache")
		return nil, nil
	}
	if err != nil {
		logger.Error("Cache failed to get user", err.Error())
	}

	var user1 user.User
	err = json.Unmarshal([]byte(userStr), &user1)
	if err != nil {
		logger.Error("Cache unmarshal user json error")
		return nil, err
	}

	endTime := time.Now()
	commons.Latency.WithLabelValues(
		commons.CacheComponent, commons.CacheGetLabel).
		Observe(float64(endTime.Sub(startTime).Milliseconds()))
	return &user1, nil
}
