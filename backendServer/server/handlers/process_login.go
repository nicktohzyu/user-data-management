package handlers

import (
	"time"
	"user-data-management/backendServer/tokenUtil"
	"user-data-management/commons"
	"user-data-management/commons/logger"
	"user-data-management/commons/packet"
	"user-data-management/commons/user"
)

func (handler handler) ProcessLogin(pkt packet.RequestPacket) packet.ResponsePacket {
	startTime := time.Now()
	Username := pkt.User.Username
	var response packet.ResponsePacket
	//TODO: password hashing
	user1, err := handler.Dbw.GetUser(Username)
	if err != nil {
		//TODO: check if error is because user does not exist
		logger.Info("Error checking password:", err)
		response = packet.ResponsePacket{
			Success: false,
			Err:     "Error",
		}
	} else if user1.Password != pkt.User.Password {
		//Check whether password is correct
		response = packet.ResponsePacket{
			Success: false,
			Err:     "Invalid password",
		}
	} else if tokenUtil.IsValidToken(user1.Token) {
		response = packet.ResponsePacket{
			Success: true,
			User: user.User{
				Username: user1.Username,
				Nickname: user1.Nickname,
				Token:    user1.Token,
			},
		}
	} else {
		//generate new token and register user
		user1.Token = tokenUtil.GenerateToken()
		err = handler.Dbw.UpdateToken(user1)
		if err != nil {
			logger.Error("Error in generating new token (login):", err)
			response = packet.ResponsePacket{
				Success: false,
				Err:     "Error in generating new token",
			}
		} else {
			response = packet.ResponsePacket{
				Success: true,
				User: user.User{
					Username: user1.Username,
					Nickname: user1.Nickname,
					Token:    user1.Token,
				},
			}
		}
	}
	endTime := time.Now()
	commons.Latency.WithLabelValues(
		commons.ServerComponent, commons.LoginLabel).
		Observe(float64(endTime.Sub(startTime).Milliseconds()))
	return response
}
