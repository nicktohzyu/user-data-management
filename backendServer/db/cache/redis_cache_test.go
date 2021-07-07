package cache

import (
	"fmt"
	"testing"
	"user-data-management/commons/logger"
	"user-data-management/commons/user"
)

var user1 = user.User{
	Username: "a",
	Password: "a",
	Nickname: "a",
}

func TestMain(m *testing.M) {
	cache := NewRedisCache()

	conn := cache.pool.Get()
	_, err := conn.Do("FLUSHDB")
	if err != nil {
		logger.Error("Error flushing redis cache")
		return
	}
	conn.Close()

	m.Run()
}

func TestAddedUser(t *testing.T) {
	cache := NewRedisCache()
	cache.Store(user1)
	fmt.Println(cache.Get(user1.Username))
}

func TestNotAddedUser(t *testing.T) {
	cache := NewRedisCache()
	ret, err := cache.Get(user1.Username)
	fmt.Println(ret, err)
}
