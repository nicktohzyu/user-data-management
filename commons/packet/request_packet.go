package packet

import "user-data-management/commons/user"

const (
	REGISTER        = "REGISTER"
	LOGIN           = "LOGIN"
	UPDATE_NICKNAME = "UPDATE_NICKNAME"
	UPDATE_IMAGE    = "UPDATE_IMAGE"
	VALIDATE        = "VALIDATE"
	LOGOUT          = "LOGOUT"
)

type RequestPacket struct {
	Format string
	User   user.User
}
