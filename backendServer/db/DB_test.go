package db

import (
	"fmt"
	"testing"
	"user-data-management/commons/user"
)

const (
	TEST_DATASOURCENAME = "root:12345678@/UDM_Testing"
)

func TestMain(m *testing.M) {
	db1 := DBWrapper{}
	db1.Init(TEST_DATASOURCENAME)
	db1.DB.Exec("DELETE FROM users")
	m.Run()
}

var user1 = user.User{
	Username: "a",
	Password: "a",
	Nickname: "a",
}

func TestDB(t *testing.T) {
	db1 := DBWrapper{}
	db1.Init(TEST_DATASOURCENAME)
}

func TestPingWithoutInit(t *testing.T) {
	//expect panic
	//TODO: pass test when panic
	db1 := DBWrapper{}
	db1.ping()
}

func TestPingAfterInit(t *testing.T) {
	db1 := DBWrapper{}
	db1.Init(TEST_DATASOURCENAME)
	db1.ping()
}

func TestNewUser(t *testing.T) {
	db1 := DBWrapper{}
	db1.Init(TEST_DATASOURCENAME)
	db1.NewUser(user1)
}

func TestUpdatePassword(t *testing.T) {
	db1 := DBWrapper{}
	db1.Init(TEST_DATASOURCENAME)
	db1.NewUser(user1)
	db1.UpdatePassword(user.User{Username: "a", Password: "b", Nickname: "b"})
}

func TestUpdateNickname(t *testing.T) {
	db1 := DBWrapper{}
	db1.Init(TEST_DATASOURCENAME)
	db1.NewUser(user1)
	db1.UpdateNickname(user.User{Username: "a", Password: "b", Nickname: "b"})
}

func TestUpdateToken(t *testing.T) {
	db1 := DBWrapper{}
	db1.Init(TEST_DATASOURCENAME)
	db1.NewUser(user1)
	db1.UpdateToken(user.User{
		Username: "a",
		Token:    string([]byte{65, 66}),
	})
}

func TestValidUser(t *testing.T) {
	db1 := DBWrapper{}
	db1.Init(TEST_DATASOURCENAME)
	db1.NewUser(user1)
	fmt.Println(db1.ValidUser(user.User{Username: "a", Password: "a"}))
	fmt.Println(db1.ValidUser(user.User{Username: "a", Password: "b"}))
}

func TestGetUser(t *testing.T) {
	db1 := DBWrapper{}
	db1.Init(TEST_DATASOURCENAME)
	db1.NewUser(user1)
	fmt.Println(db1.GetUser("a"))
}

func TestGetUserTwice(t *testing.T) {
	//for testing cache
	db1 := DBWrapper{}
	db1.Init(TEST_DATASOURCENAME)
	db1.NewUser(user1)
	db1.GetUser("a")
	db1.GetUser("a")
}
