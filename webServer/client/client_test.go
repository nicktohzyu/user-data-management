package client

import (
	"fmt"
	"testing"
	"user-data-management/commons/user"
)

var user1 = user.User{
	Username: "id1",
	Password: "pass1",
	Nickname: "a",
	Token:    "xxx xxx xxx xxx ",
}

func TestClient_login(t *testing.T) {
	client, _ := InitClient(":8001")
	rsp, err := client.Login(user1)
	fmt.Println(rsp, err)
}

func TestClient_register(t *testing.T) {
	client, _ := InitClient(":8001")
	rsp, err := client.register(user1)
	fmt.Println(rsp, err)
}

func TestClient_updateNickname(t *testing.T) {
	client, _ := InitClient(":8001")
	rsp, err := client.updateNickname(user.User{
		Username: user1.Username,
		Token:    user1.Token,
		Nickname: "C",
	})
	fmt.Println(rsp, err)
}

func TestClient_logout(t *testing.T) {
	client, _ := InitClient(":8001")
	rsp, err := client.logout(user.User{
		Username: user1.Username,
		Token:    user1.Token,
	})
	fmt.Println(rsp, err)
}

func TestClient_requestWhenLoggedOut(t *testing.T) {
	client, _ := InitClient(":8001")
	rsp, err := client.updateNickname(user.User{
		Username: user1.Username,
		Token:    "0000000000000000",
		Nickname: "C",
	})
	fmt.Println(rsp, err)
}

func TestClient_register_login_updateNickname(t *testing.T) {
	client, _ := InitClient(":8001")
	rsp, err := client.register(user1)
	fmt.Println(rsp, err)
	rsp, err = client.Login(user1)
	fmt.Println(rsp, err)
	fmt.Println("Received token: ", []byte(rsp.Token))
	rsp, err = client.updateNickname(user.User{
		Username: rsp.Username,
		Token:    rsp.Token,
		Nickname: "b",
	})
	fmt.Println(rsp, err)
}
