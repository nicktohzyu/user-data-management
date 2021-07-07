package main

import (
	crand "crypto/rand"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/time/rate"
	"math/big"
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
	NUM_ITERATIONS          = 10000
	REQUESTS_PER_ITERATION  = 600
	MAX_USERID              = 100*10000 - 1
	TOKEN                   = "xxx xxx xxx xxx "
	NUM_USERS               = 1000
	NUM_WORKERS             = 500
	TIME_BETWEEN_ITERATIONS = time.Millisecond * 100
	TIME_BETWEEN_REQUESTS   = 0
	LIMITER_QPS             = 4000
	LIMITER_BURST           = 500
)

var clientPtr *client.Client
var userChans [NUM_WORKERS]chan user.User
var resultsChan chan string = make(chan string, REQUESTS_PER_ITERATION)

func main() {
	//test data should be inserted into DB, and backend server should be listening
	prometheus.MustRegister(commons.Latency)
	http.Handle("/client_metrics", promhttp.Handler())
	var err error
	clientPtr, err = client.InitClient(":8001")
	if err != nil {
		panic(err)
	}
	const usersPerWorker = (REQUESTS_PER_ITERATION + NUM_WORKERS - 1) / NUM_WORKERS //divide ceil
	for i := 0; i < NUM_WORKERS; i++ {
		userChans[i] = make(chan user.User, usersPerWorker)
		go worker(userChans[i], resultsChan)
	}
	go func() {
		//	rate limiter
		users := generateUsers()
		for i := 0; i < NUM_ITERATIONS; i++ {
			fmt.Println("Iteration ", i)
			go iteration(users)
			time.Sleep(TIME_BETWEEN_ITERATIONS)
		}
	}()

	http.ListenAndServe(":8081", nil)
	logger.Info("Prometheus server started")
}

func iteration(users *[]user.User) {
	startTime := time.Now()
	for i := 0; i < REQUESTS_PER_ITERATION; i++ {
		user1 := (*users)[rand.Int()%NUM_USERS]
		go func() {
			userChans[i%NUM_WORKERS] <- user1
		}()
		time.Sleep(TIME_BETWEEN_REQUESTS)
	}
	results := make(map[string]int)
	for i := 0; i < REQUESTS_PER_ITERATION; i++ {
		result := <-resultsChan
		results[result]++
	}
	fmt.Println(results)
	endTime := time.Now()
	fmt.Println("time taken: ", endTime.Sub(startTime))
}

func generateUsers() *[]user.User {
	users := make([]user.User, NUM_USERS)
	nums := generateNumbers()
	for i := 0; i < NUM_USERS; i++ {
		num := (*nums)[i]
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

//generateNumbers returns a pointer to an array of unique numbers
func generateNumbers() *[]int {
	nums := make([]int, NUM_USERS)
	for i := 0; i < NUM_USERS; i++ {
		used := true
		var num int
		for used {
			used = false
			bigNum, _ := crand.Int(crand.Reader, big.NewInt(int64(MAX_USERID)))
			num = int(bigNum.Int64())
			for _, v := range nums {
				if num == v {
					used = true
				}
			}
		}
		nums[i] = num
	}
	return &nums
}

var limiter = rate.NewLimiter(LIMITER_QPS, LIMITER_BURST)

func worker(userChan chan user.User, res chan string) {
	for user1 := range userChan {
		if !limiter.Allow() {
			res <- "dropped by limiter"
			continue
		}
		startTime := time.Now()
		resp, err := clientPtr.Login(user1)
		_ = resp //no use for resp for now
		endTime := time.Now()
		elapsedTime := endTime.Sub(startTime)
		commons.Latency.WithLabelValues(
			commons.StressTestComponent, commons.RoundTripLabel).
			Observe(float64(elapsedTime.Milliseconds()))
		if err != nil {
			res <- err.Error()
		}
		res <- "success"
	}
}
