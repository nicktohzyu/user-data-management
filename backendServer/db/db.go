package db

import (
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"time"
	"user-data-management/backendServer/db/cache"
	"user-data-management/commons"
	"user-data-management/commons/logger"
	"user-data-management/commons/user"
)

type DBWrapper struct {
	DB     *sql.DB
	IsInit bool
	cache  *cache.RedisCache
}

func (dbw *DBWrapper) Init(dataSourceName string) {
	loggerLevel := logger.Level
	logger.Level = logger.INFO
	logger.Info("Initializing db wrapper")
	defer func() { logger.Level = loggerLevel }()

	if dbw.IsInit {
		panic("Already initialized")
	}
	var err error
	dbw.DB, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		panic(err.Error())
	}
	dbw.DB.SetConnMaxLifetime(time.Second * DBMaxConnectionLifetime)
	dbw.DB.SetMaxIdleConns(DBMaxIdleConnections)
	dbw.DB.SetMaxOpenConns(DBMaxOpenConnections)
	//Validate DSN data:
	err = dbw.DB.Ping()
	if err != nil {
		panic(err.Error())
	}
	dbw.IsInit = true

	//init cache
	dbw.cache = cache.NewRedisCache()

	logger.Info("DB and cache initialized")
}

func (dbw *DBWrapper) ping() {
	if !dbw.IsInit {
		panic("DB not initialized")
	}
	var err error
	err = dbw.DB.Ping()
	if err != nil {
		panic(err.Error())
	}
}

func (dbw *DBWrapper) NewUser(user user.User) error {
	db := dbw.DB
	_, err := db.Query("INSERT INTO users VALUES (?, ?, ?, ?)",
		user.Username,
		user.Password,
		user.Nickname,
		user.Token,
	)
	if err != nil {
		logger.Error("db insertion error: " + err.Error())
		return errors.New("db insertion error: " + err.Error())
	}

	// add user in cache
	err = dbw.cache.Store(user)
	if err != nil {
		logger.Error("NewUser: error storing user into cache", err)
		//no need to return err here as long as DB has correct data
	}
	return nil
}

func (dbw *DBWrapper) IsExists(username string) (bool, error) {
	db := dbw.DB
	var exists bool
	err := db.QueryRow(
		"SELECT exists (SELECT * FROM users WHERE username = ?)",
		username).
		Scan(&exists)
	if err != nil {
		panic("db query error: " + err.Error())
		return false, err
	}
	if !exists {
		return false, nil
	}
	return true, nil
}

func (dbw *DBWrapper) UpdatePassword(user user.User) error {
	db := dbw.DB
	//check if user exists
	exists, err := dbw.IsExists(user.Username)
	if err != nil {
		//	TODO: handle
	}
	if !exists {
		return errors.New("user not found")
	}
	//update
	_, err = db.Query(`UPDATE users set password = ?
		where username = ?`,
		user.Password,
		user.Username,
	)
	if err != nil {
		panic("db update error: " + err.Error())
	}

	// add/update user in cache
	err = dbw.cache.Store(user)
	if err != nil {
		logger.Error("UpdatePassword: error storing user into cache", err)
		return err //must return err because cache may have incorrect info
	}
	return nil
}

func (dbw *DBWrapper) UpdateNickname(user user.User) error {
	db := dbw.DB
	//update
	_, err := db.Query(`UPDATE users set nickname = ?
		where username = ?`,
		user.Nickname,
		user.Username,
	)
	if err != nil {
		return err
	}

	// add/update user in cache
	err = dbw.cache.Store(user)
	if err != nil {
		logger.Error("UpdateNickname: error storing user into cache", err)
		return err //must return err because cache may have incorrect info
	}
	return nil
}

func (dbw *DBWrapper) UpdateToken(user user.User) error {
	startTime := time.Now()
	db := dbw.DB
	//update
	res, err := db.Query(`UPDATE users set token = ?
		where username = ?`,
		user.Token,
		user.Username,
	)
	res.Close()
	if err != nil {
		logger.Error("UpdateToken DB error: ", err)
		return err
	}

	// add/update user in cache
	err = dbw.cache.Store(user)
	if err != nil {
		logger.Error("UpdateToken: error storing user into cache", err)
		return err //must return err because cache may have incorrect info
	}

	endTime := time.Now()
	commons.Latency.WithLabelValues(
		commons.DBComponent, commons.UpdateTokenLabel).
		Observe(float64(endTime.Sub(startTime).Milliseconds()))
	return nil
}

// ValidUser tests if username and password match
func (dbw *DBWrapper) ValidUser(user user.User) bool {
	db := dbw.DB
	//check if user exists
	var match bool
	err := db.QueryRow(
		"SELECT exists (SELECT * FROM users WHERE username = ? and password = ?)",
		user.Username,
		user.Password,
	).
		Scan(&match)
	if err != nil {
		panic("db query error: " + err.Error())
	}
	return match
}

func (dbw *DBWrapper) GetUser(queryUsername string) (user.User, error) {
	loggerLevel := logger.Level
	logger.Level = logger.ERROR
	defer func() { logger.Level = loggerLevel }()

	startTime := time.Now()
	var user1 user.User
	// check cache
	ret, err := dbw.cache.Get(queryUsername)
	if err != nil {
		logger.Error("Cache retrieval error: ", err)
		return user.User{}, err
	} else if ret != nil {
		logger.Info("GetUser: user found in cache")
		user1 = *ret

		endTime := time.Now()
		commons.Latency.WithLabelValues(
			commons.DBComponent, commons.GetUserCacheLabel).
			Observe(float64(endTime.Sub(startTime).Milliseconds()))
		return user1, nil
	}
	//not in cache, check DB then add to cache
	logger.Info("GetUser: user not found in cache")
	db := dbw.DB
	var username1, password1, nickname1, token string
	rows, err := db.Query(
		"SELECT * FROM users WHERE username = ?",
		queryUsername)
	if err != nil {
		logger.Error("DB query error (login): ", err)
		return user.User{}, err
	}
	rows.Next()
	err = rows.Scan(&username1, &password1, &nickname1, &token)
	if err != nil {
		logger.Error("DB result scan error (login): ", err)
		return user.User{}, err
	}
	err = rows.Close()
	if err != nil {
		logger.Error("DB close rows error (login): ", err)
		return user.User{}, err
	}
	user1 = user.User{
		Username: username1,
		Password: password1,
		Nickname: nickname1,
		Token:    token,
	}
	if err != nil {
		//TODO: differentiate between user not found vs true db error
		logger.Error("db query error: ", err)
		return user.User{}, err
	}
	//	add user to cache
	err = dbw.cache.Store(user1)
	//TODO: invalidate/update cache if user properties modified
	if err != nil {
		logger.Error("Error storing user into cache", err)
	}

	endTime := time.Now()
	commons.Latency.WithLabelValues(
		commons.DBComponent, commons.GetUserDBLabel).
		Observe(float64(endTime.Sub(startTime).Milliseconds()))
	return user1, nil
}
