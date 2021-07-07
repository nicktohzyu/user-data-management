package packet

import "user-data-management/commons/user"

type ResponsePacket struct {
	Success bool
	User    user.User
	Err     string
}
