package client

import (
	"encoding/json"
	"errors"
	"time"
	"user-data-management/commons/logger"
	"user-data-management/commons/packet"
	"user-data-management/commons/user"
	"user-data-management/webServer/client/pool"
)

type Client struct {
	pool *pool.ConnectionPool
	//	can have many outgoing TCP connections
}

func InitClient(address string) (*Client, error) {

	poolConfig := pool.Config{
		NumConnections: pool.NUM_CONNECTIONS,
		RetryDuration:  time.Millisecond * 100,
		WaitDuration:   time.Millisecond * 100,
		ReadDuration:   time.Millisecond * 100,
	}
	connPool := pool.NewPool(poolConfig, address)

	//c1, err := createConnection(address)
	//if err != nil {
	//	fmt.Println(err)
	//	return nil, err
	//}
	client := &Client{
		pool: connPool,
	}
	return client, nil
}

func (client Client) sendPacket(outPkt packet.RequestPacket) (*packet.ResponsePacket, error) {

	outBytes, err := json.Marshal(outPkt)
	if err != nil {
		logger.Error("Error marshalling out packet (login):", err)
		return nil, err
	}

	//use pool to send the packet
	response, err := client.pool.SendBytes(outBytes)
	if err != nil {
		//TODO: handle
		return nil, err
	}

	// unmarshal the packet
	var inPkt packet.ResponsePacket
	err = json.Unmarshal(*response, &inPkt)
	if err != nil {
		//TODO: handle
		return nil, err
	}

	if !inPkt.Success {
		logger.Error("request unsuccessful:", inPkt.Err)
		return &inPkt, errors.New(inPkt.Err)
	}
	return &inPkt, nil
}

func (client Client) register(user1 user.User) (*user.User, error) {
	outPkt := packet.RequestPacket{
		Format: packet.REGISTER,
		User: user.User{
			Username: user1.Username,
			Password: user1.Password,
			Nickname: user1.Nickname,
		},
	}

	inPktPtr, err := client.sendPacket(outPkt)
	if err != nil {
		//TODO: handle
		return nil, err
	}

	inPkt := *inPktPtr
	//process response
	if !inPkt.Success {
		//TODO: handle
		return nil, err
	}
	return &user.User{
		Username: inPkt.User.Username,
		Token:    inPkt.User.Token,
		Nickname: inPkt.User.Nickname,
	}, err
}

func (client Client) Login(user1 user.User) (*user.User, error) {
	outPkt := packet.RequestPacket{
		Format: packet.LOGIN,
		User: user.User{
			Username: user1.Username,
			Password: user1.Password,
		},
	}

	inPktPtr, err := client.sendPacket(outPkt)
	if err != nil {
		//TODO: handle
		return nil, err
	}

	inPkt := *inPktPtr
	//process response
	if !inPkt.Success {
		//TODO: handle
		return nil, err
	}
	return &user.User{
		Username: inPkt.User.Username,
		Token:    inPkt.User.Token,
		Nickname: inPkt.User.Nickname,
	}, err
}

func (client Client) logout(user1 user.User) (*user.User, error) {
	outPkt := packet.RequestPacket{
		Format: packet.LOGOUT,
		User: user.User{
			Username: user1.Username,
			Token:    user1.Token,
		},
	}

	inPktPtr, err := client.sendPacket(outPkt)
	if err != nil {
		//TODO: handle
		return nil, err
	}

	inPkt := *inPktPtr
	//process response
	if !inPkt.Success {
		//TODO: handle
		return nil, err
	}
	return &user.User{
		Username: inPkt.User.Username,
	}, err
}

func (client Client) updateNickname(user1 user.User) (*user.User, error) {
	outPkt := packet.RequestPacket{
		Format: packet.UPDATE_NICKNAME,
		User: user.User{
			Username: user1.Username,
			Token:    user1.Token,
			Nickname: user1.Nickname,
		},
	}

	inPktPtr, err := client.sendPacket(outPkt)
	if err != nil {
		//TODO: handle
		return nil, err
	}

	inPkt := *inPktPtr
	//process response
	if !inPkt.Success {
		//TODO: handle
		return nil, err
	}
	return &user.User{
		Username: inPkt.User.Username,
		Nickname: inPkt.User.Nickname,
	}, err
}
