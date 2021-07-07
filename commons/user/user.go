package user

type User struct {
	Username string
	//TODO: store password as hash instead
	Password string
	Nickname string
	Token    string // 16 bytes
}
