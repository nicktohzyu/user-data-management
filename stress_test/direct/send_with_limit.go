package main

import (
	"golang.org/x/time/rate"
	"time"
	"user-data-management/commons"
	"user-data-management/commons/user"
	"user-data-management/webServer/client"
)

var limiter = rate.NewLimiter(REQUESTS_PER_ITERATION, REQUESTS_PER_ITERATION)

func sendWithLimit(client *client.Client, user1 user.User, iterationNum int) string {
	if !limiter.Allow() {
		return "dropped by limiter"
	}
	startTime := time.Now()
	resp, err := (*client).Login(user1)
	_ = resp //no use for resp for now
	endTime := time.Now()
	elapsedTime := endTime.Sub(startTime)
	commons.Latency.WithLabelValues(
		commons.StressTestComponent, commons.RoundTripLabel).
		Observe(float64(elapsedTime))
	if err != nil {
		return err.Error()
	}
	return "success"
}
