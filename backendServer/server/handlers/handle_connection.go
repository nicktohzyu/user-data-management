package handlers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"user-data-management/commons/logger"
	"user-data-management/commons/packet"
)

func (handler handler) HandleConnection(c net.Conn) {
	logger.Info("Serving %s\n", c.RemoteAddr().String())
	for {
		logger.Debug("Listening")
		inStr, err := bufio.NewReader(c).ReadBytes('\n')
		if err != nil {
			if err.Error() != "EOF" {
				logger.Error("read error:", err)
			}
			break
		}
		logger.Debug("Backend received string:", string(inStr))
		if strings.TrimSpace(string(inStr)) == "STOP" {
			logger.Debug("Closing connection")
			break
		}

		outPktStr, err := handler.handlePacketBytes(&inStr)
		if err != nil {
			logger.Error("Error handling packet:", err)
			//TODO: send response indicating server error
			continue
		}

		logger.Debug("Sending response: ", string(*outPktStr)+"\n")
		_, err = fmt.Fprintf(c, string(*outPktStr)+"\n")
		if err != nil {
			logger.Error("send response error:", err)
		}
		logger.Debug("Response sent")
	}
	c.Close()
}

//invalid input strings should not cause this function to return an error
func (handler handler) handlePacketBytes(inBytes *[]byte) (*[]byte, error) {
	var (
		inPkt  packet.RequestPacket
		outPkt packet.ResponsePacket
	)
	err := json.Unmarshal(*inBytes, &inPkt)
	if err != nil {
		logger.Error("unmarshal error:", err)
		outPkt = packet.ResponsePacket{
			Success: false,
			Err:     "Invalid packet",
		}
	} else {
		outPkt = handler.ProcessPacket(inPkt)
	}
	outBytes, err := json.Marshal(outPkt)
	if err != nil {
		logger.Error("Marshal outgoing packet error")
	}
	return &outBytes, err
}
