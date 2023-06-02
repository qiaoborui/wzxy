package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"
	"wobuzaixiaoyuan/utils"
	"wobuzaixiaoyuan/wzxy"
)

// TEST GITHUB ACTION
func main() {
	log.SetFlags(log.Ltime | log.Ldate)
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		err = os.Mkdir("logs", 0755)
		if err != nil {
			log.Fatal(err)
		}
	}
	dateNow := getDate()
	dateTmp := ""
	timeTmp := time.Now()
	storage, err := utils.FetchData("{\"status\": 1 }")
	eventMap := make(map[string]int)
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		timeNow := time.Now()
		if timeNow.Sub(timeTmp).Minutes() >= 30 {
			timeTmp = timeNow
			storage, err = utils.FetchData("{\"status\": 1 }")
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		if dateNow != dateTmp {
			dateTmp = dateNow
			for _, user := range storage.Results {
				eventMap[user.RealName] = 0
			}
		}
		var users []*wzxy.User
		results := make(chan string, 10)
		for _, user := range storage.Results {
			if timeNow.Format("2260") > user.Start && timeNow.Format("2260") < user.End && eventMap[user.RealName] < 2 {
				users = append(users, &wzxy.User{
					RealName: user.RealName,
					Username: user.Username,
					Password: user.Password,
					Result:   results,
				})
				eventMap[user.RealName]++
			}
		}
		if len(users) == 0 {
			log.SetOutput(os.Stdout)
			log.Printf("没有需要打卡的用户")
		}
		if len(users) != 0 {
			doWork(users)
		}
		time.Sleep(5 * time.Minute)
	}
}
func doWork(users []*wzxy.User) {
	timeNow := time.Now()
	logFileName := fmt.Sprintf("logs/%d-%02d-%02d %02d.%02d.%02d.log", timeNow.Year(), timeNow.Month(), timeNow.Day(), timeNow.Hour(), timeNow.Minute(), timeNow.Second())
	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0777)
	multiWriter := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(multiWriter)
	if err != nil {
		log.Fatal(err)
	}
	var wg sync.WaitGroup
	start := time.Now()

	maxConcurrent := 10

	// 用来控制最大并发数量
	sem := make(chan struct{}, maxConcurrent)

	// 用来接收已签到的wzxy对象
	successCh := make(chan *wzxy.Session, len(users))

	//var wg sync.WaitGroup

	for _, user := range users {
		w := wzxy.NewWzxy(user)
		sem <- struct{}{}
		wg.Add(1)
		go func(w *wzxy.Session) {
			defer func() { <-sem }()
			defer wg.Done()

			if err := w.Login(); err != nil {
				fmt.Printf("[%s] login failed: %v\n", w.User.RealName, err)
				return
			}

			_ = w.Sign()

			successCh <- w
		}(w)
	}

	wg.Wait()
	close(sem)
	close(successCh)

	// 按顺序输出已签到的wzxy对象
	for w := range successCh {
		log.Printf("%s\n", <-w.User.Result)
	}
	elapsed := time.Since(start)
	log.Printf("程序运行时间为：%s \n", elapsed)
	logFile.Close()
}
func getDate() string {
	return time.Now().Format("20060102")
}
