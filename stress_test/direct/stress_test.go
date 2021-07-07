package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"math/rand"
	"net/http"
	"strconv"
	"time"
	"user-data-management/commons"
	"user-data-management/commons/logger"
	"user-data-management/commons/user"
	"user-data-management/webServer/client"
)

const (
	NUM_ITERATIONS         = 10
	REQUESTS_PER_ITERATION = 100
	MAX_USERID             = 2 * 10000
	TOKEN                  = "xxx xxx xxx xxx "
	NUM_USERS              = 200
)

func main() {
	//test data should be inserted into DB, and backend server should be listening
	prometheus.MustRegister(commons.Latency)
	http.Handle("/client_metrics", promhttp.Handler())
	go func() {
		client, err := client.InitClient(":8001")
		if err != nil {
			panic(err)
		}
		//	rate limiter
		users := generateUsers()
		for i := 0; i < NUM_ITERATIONS; i++ {
			fmt.Println("Iteration ", i)
			iteration(client, users, i)
			time.Sleep(time.Millisecond * 200)
		}
	}()

	http.ListenAndServe(":8081", nil)
	logger.Info("Prometheus server started")
}

func iteration(client *client.Client, users *[]user.User, iterationNum int) {
	startTime := time.Now()
	resultChannel := make(chan string, REQUESTS_PER_ITERATION)
	for i := 0; i < REQUESTS_PER_ITERATION; i++ {
		user1 := (*users)[rand.Int()%NUM_USERS]
		go func() {
			resultChannel <- sendWithLimit(client, user1, iterationNum)
		}()
		time.Sleep(time.Millisecond / 2)
	}
	results := make(map[string]int)
	for i := 0; i < REQUESTS_PER_ITERATION; i++ {
		result := <-resultChannel
		results[result]++
	}
	fmt.Println(results)
	endTime := time.Now()
	fmt.Println("time taken: ", endTime.Sub(startTime))
}

func generateUsers() *[]user.User {
	users := make([]user.User, NUM_USERS)
	for i := 0; i < NUM_USERS; i++ {
		num := rand.Int() % MAX_USERID
		id := "id" + strconv.Itoa(num)
		nickname := "nickname" + strconv.Itoa(num%10000)
		pass := "pass" + strconv.Itoa(num%10000)
		user1 := user.User{
			Username: id,
			Nickname: nickname,
			Password: pass,
			Token:    TOKEN,
		}
		users[i] = user1
	}
	return &users
}
