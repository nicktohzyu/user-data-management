package handlers

import (
	"user-data-management/backendServer/tokenUtil"
	"user-data-management/commons/logger"
	"user-data-management/commons/packet"
	"user-data-management/commons/user"
)

func (handler handler) ProcessLogout(pkt packet.RequestPacket) packet.ResponsePacket {
	//retrieve user data
	user1, err := handler.Dbw.GetUser(pkt.User.Username)
	if err != nil {
		logger.Error("Error retreiving user data:", err)
		return packet.ResponsePacket{
			Success: false,
			Err:     "Error",
		}
	}

	//check if token is correct
	if user1.Token != pkt.User.Token {
		return packet.ResponsePacket{
			Success: false,
			Err:     "Invalid token",
		}
	}

	//generate new token and register user
	user1.Token = tokenUtil.SPECIAL
	err = handler.Dbw.UpdateToken(user1)
	if err != nil {
		logger.Error("Error updating token (logout):", err)
		return packet.ResponsePacket{
			Success: false,
			Err:     "Error updating token",
		}
	}
	return packet.ResponsePacket{
		Success: true,
		User: user.User{
			Username: user1.Username,
		},
	}
}
