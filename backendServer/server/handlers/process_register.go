package handlers

import (
	"user-data-management/backendServer/tokenUtil"
	"user-data-management/commons/logger"
	"user-data-management/commons/packet"
	"user-data-management/commons/user"
)

func (handler handler) ProcessRegister(pkt packet.RequestPacket) packet.ResponsePacket {
	user1 := pkt.User

	//check whether username is already registered
	exists, err := handler.Dbw.IsExists(user1.Username)
	if exists {
		return packet.ResponsePacket{
			Success: false,
			Err:     "User already exists",
		}
	}

	//TODO: check that username is within max length

	//generate new token and register user
	user1.Token = tokenUtil.GenerateToken()
	err = handler.Dbw.NewUser(user1)
	if err != nil {
		logger.Error("process register:", err)
		return packet.ResponsePacket{
			Success: false,
			Err:     "Error in registration",
		}
	}
	return packet.ResponsePacket{
		Success: true,
		User: user.User{
			Username: user1.Username,
			Nickname: user1.Nickname,
			Token:    user1.Token,
		},
	}
}
