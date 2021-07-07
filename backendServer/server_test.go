package main

import (
	"encoding/json"
	"fmt"
	"testing"
	"user-data-management/backendServer/db"
	"user-data-management/backendServer/server/handlers"
	"user-data-management/commons/packet"
	"user-data-management/commons/user"
)

const(
	dataSourceName = "root:12345678@/UDM_Testing"
)

func TestMain(m *testing.M) {
	db1 := db.DBWrapper{}
	db1.Init("root:12345678@/UDM_Testing")
	db1.DB.Exec("DELETE FROM users")
	m.Run()
}

var user1 = user.User{
	Username: "a",
	Password: "a",
	Nickname: "a",
	Token:    "AAA AAA AAA AAA ",
}

func TestRegister(t *testing.T) {
	pkt := packet.RequestPacket{
		Format: packet.REGISTER,
		User:   user1,
	}
	handler := handlers.NewHandler(dataSourceName)
	rsp := handler.ProcessRegister(pkt)
	_ = rsp
}

func TestLogin(t *testing.T) {
	pkt := packet.RequestPacket{
		Format: packet.LOGIN,
		User:   user1,
	}
	handler := handlers.NewHandler(dataSourceName)
	handler.Dbw.NewUser(user1)
	rsp := handler.ProcessLogin(pkt)
	str, _ := json.Marshal(rsp)
	fmt.Println(string(str))
}

func TestUpdateNickname(t *testing.T) {
	handler := handlers.NewHandler(dataSourceName)
	handler.Dbw.NewUser(user1)
	handler.Dbw.UpdateToken(user1)

	pkt := packet.RequestPacket{
		Format: packet.UPDATE_NICKNAME,
		User: user.User{
			Username: user1.Username,
			Token:    user1.Token,
			Nickname: "B",
		},
	}
	rsp := handler.ProcessUpdateNickname(pkt)
	str, _ := json.Marshal(rsp)
	fmt.Println(string(str))
}

func TestLogout(t *testing.T) {
	handler := handlers.NewHandler(dataSourceName)
	handler.Dbw.NewUser(user1)
	handler.Dbw.UpdateToken(user1)

	pkt := packet.RequestPacket{
		Format: packet.LOGOUT,
		User: user.User{
			Username: user1.Username,
			Token:    user1.Token,
		},
	}
	rsp := handler.ProcessLogout(pkt)
	str, _ := json.Marshal(rsp)
	fmt.Println(string(str))
}

func TestProcessPacketRegister(t *testing.T) {
	pkt := packet.RequestPacket{
		Format: packet.REGISTER,
		User:   user1,
	}
	handler := handlers.NewHandler(dataSourceName)
	handler.ProcessPacket(pkt)
}

func TestRegisterTCP(t *testing.T) {
	pkt := packet.RequestPacket{
		Format: packet.REGISTER,
		User:   user1,
	}
	pktStr, _ := json.Marshal(pkt)
	fmt.Println(string(pktStr))
	//{"Format":"REGISTER","User":{"Username":"a","Password":"a","Nickname":"a","Token":""}}
}
