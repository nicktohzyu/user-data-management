package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"os"
	"strconv"
	"user-data-management/backendServer/server"
	"user-data-management/commons"
	"user-data-management/commons/logger"
)

const (
	SERVER_PORT     = 8001
	PROMETHEUS_PORT = 8080
)

func main() {
	DATASOURCENAME := os.Getenv("DSN")
	//"root:12345678@/UDM_Testing"
	//"root:12345678@/UDM_StressTest"
	prometheus.MustRegister(commons.Latency)
	http.Handle("/db_metrics", promhttp.Handler())
	go http.ListenAndServe(":"+strconv.Itoa(PROMETHEUS_PORT), nil)
	logger.Info("Prometheus server started")

	server1 := server.InitServer(SERVER_PORT, DATASOURCENAME)
	server1.HandleConnections()
}
