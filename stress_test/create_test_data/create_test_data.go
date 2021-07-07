package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strings"
	"time"
)

const TOKEN = "xxx xxx xxx xxx "

func main() {
	startTime := time.Now()
	db, err := sql.Open("mysql", "root:12345678@/UDM_StressTest")
	if err != nil {
		fmt.Println("Create test data script: Unable to open the DB")
		return
	}
	defer db.Close()
	db.SetConnMaxLifetime(time.Second * 6000)
	db.SetMaxIdleConns(2048)
	db.SetMaxOpenConns(2048)
	db.Exec("DELETE FROM users")
	for j := 0; j < 100; j++ {
		vals := make([]string, 10000)
		for i := 0; i < 10000; i++ {
			id := fmt.Sprintf("id%d", j*10000+i)
			pass := fmt.Sprintf("pass%d", i)
			nick := fmt.Sprintf("nick%d", i)

			curr := fmt.Sprintf("('%s', '%s', '%s', '%s')", id, pass, nick, TOKEN)
			vals[i] = curr
		}

		queryStr := fmt.Sprintf("INSERT INTO users(username, password, nickname, token) VALUES %s", strings.Join(vals, ","))

		insert, err := db.Query(queryStr)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		insert.Close()

		fmt.Println("Done iteration: %v", j)

		time.Sleep(100 * time.Millisecond)
	}
	endTime := time.Now()
	fmt.Println("Time taken: ", endTime.Sub(startTime))
}
