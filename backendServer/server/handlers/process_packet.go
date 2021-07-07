package handlers

import (
	"time"
	"user-data-management/commons"
	"user-data-management/commons/logger"
	"user-data-management/commons/packet"
)

func (handler handler) ProcessPacket(pkt packet.RequestPacket) packet.ResponsePacket {
	startTime := time.Now()
	var returnPkt packet.ResponsePacket
	switch pkt.Format {
	case packet.REGISTER:
		returnPkt = handler.ProcessRegister(pkt)
	case packet.LOGIN:
		returnPkt = handler.ProcessLogin(pkt)
	case packet.UPDATE_NICKNAME:
		returnPkt = handler.ProcessUpdateNickname(pkt)
	case packet.UPDATE_IMAGE:
	case packet.VALIDATE:
	case packet.LOGOUT:
		returnPkt = handler.ProcessLogout(pkt)
	default:
		//TODO: reply with invalid request format error message
		logger.Error("invalid request packet format")
		returnPkt = packet.ResponsePacket{Success: false}
	}
	endTime := time.Now()
	commons.Latency.WithLabelValues(
		commons.ServerComponent, commons.ProcessPacketLabel).
		Observe(float64(endTime.Sub(startTime).Milliseconds()))
	return returnPkt
}
