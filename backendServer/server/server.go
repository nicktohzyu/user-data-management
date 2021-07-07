package server

import (
	"net"
	"strconv"
	"user-data-management/backendServer/server/handlers"
	"user-data-management/commons/logger"
)

type Server struct {
	port           int
	dataSourceName string
}

func InitServer(port int, dataSourceName string) Server {
	return Server{
		port:           port,
		dataSourceName: dataSourceName,
	}
}

func (server Server) HandleConnections() {
	l, err := net.Listen("tcp4", ":"+strconv.Itoa(server.port))
	if err != nil {
		logger.Error("Error opening port for connections:", err)
		return
	}
	defer l.Close()
	handler1 := handlers.NewHandler(server.dataSourceName)
	for {
		c, err := l.Accept()
		if err != nil {
			logger.Error("Error accepting connection:", err)
			return
		}
		go handler1.HandleConnection(c)
	}
}
