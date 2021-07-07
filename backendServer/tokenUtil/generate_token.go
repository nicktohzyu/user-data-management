package tokenUtil

import (
	"crypto/rand"
	"math/big"
	"time"
	"user-data-management/commons"
	"user-data-management/commons/logger"
)

const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func GenerateToken() string {
	startTime := time.Now()
	arr := make([]byte, LENGTH)
	for i, _ := range arr {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		arr[i] = letters[num.Int64()]
	}
	str := string(arr)
	endTime := time.Now()
	commons.Latency.WithLabelValues(
		commons.ServerComponent, commons.TokenLabel).
		Observe(float64(endTime.Sub(startTime).Milliseconds()))
	if ValidTokenFormat(str) {
		return str
	}
	logger.Info("invalid token generated: ", str)
	return GenerateToken()
}
