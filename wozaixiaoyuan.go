package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
	"wobuzaixiaoyuan/utils"
	"wobuzaixiaoyuan/wzxy"
)

// TEST GITHUB ACTION
func main() {
	log.SetFlags(log.Ltime | log.Ldate)
	startLogServer()
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
			if !CompareTime(user.Start) && CompareTime(user.End) && eventMap[user.RealName] < 2 {
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
	_ = logFile.Close()
}
func getDate() string {
	return time.Now().Format("20060102")
}
func CompareTime(inputTime string) bool {
	//构造包含当前日期的字符串
	dateStr := time.Now().Format("2006-01-02") //使用Go语言规定的"2006-01-02"作为日期格式

	//将输入时间和日期信息组合成一个完整的时间字符串
	fullTimeStr := fmt.Sprintf("%s %s", dateStr, inputTime)

	//解析时间字符串为time类型，获取输入时间
	layout := "2006-01-02 15:04"
	t, err := time.ParseInLocation(layout, fullTimeStr, time.Local)
	if err != nil {
		fmt.Println(err)
		return false
	}

	//获取当前时间
	currentTime := time.Now()

	//比较当前时间和输入时间的差值
	if t.After(currentTime) {
		//fmt.Printf("%s is after current time %s", inputTime, currentTime.Format(layout))
		return true
	} else if t.Equal(currentTime) {
		//fmt.Printf("%s is equal to current time %s", inputTime, currentTime.Format(layout))
		return false
	} else {
		//fmt.Printf("%s is before current time %s", inputTime, currentTime.Format(layout))
		return false
	}
}
func startLogServer() {
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		err = os.Mkdir("logs", 0755)
		if err != nil {
			log.Fatal(err)
		}
	}
	port := os.Getenv("PORT") // 读取环境变量中的端口号
	if port == "" {
		port = "8000" // 默认端口号为8000
	}

	go func() {
		fs := http.FileServer(http.Dir("./logs"))
		// 监听指定的端口号
		log.Printf("Listening on :%s...\n", port)
		_ = http.ListenAndServe(":"+port, fs)
	}()
}
