package handlers

import (
	"user-data-management/backendServer/tokenUtil"
	"user-data-management/commons/logger"
	"user-data-management/commons/packet"
	"user-data-management/commons/user"
)

func (handler handler) ProcessUpdateNickname(pkt packet.RequestPacket) packet.ResponsePacket {
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
	if user1.Token != pkt.User.Token || pkt.User.Token == tokenUtil.SPECIAL {
		return packet.ResponsePacket{
			Success: false,
			Err:     "Invalid token",
		}
	}

	//store new nickname
	user1.Nickname = pkt.User.Nickname
	err = handler.Dbw.UpdateNickname(user1)
	if err != nil {
		logger.Error("Error in storing new nickname:", err)
		return packet.ResponsePacket{
			Success: false,
			Err:     "Error",
		}
	}
	return packet.ResponsePacket{
		Success: true,
		User: user.User{
			Username: user1.Username,
			Nickname: user1.Nickname,
		},
	}
}
