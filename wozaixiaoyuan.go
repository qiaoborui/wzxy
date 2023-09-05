package main

import (
	"bytes"
	"fmt"
	"github.com/robfig/cron"
	"io"
	"log"
	"os"
	"sync"
	"time"
	"wobuzaixiaoyuan/pkg/common"
	"wobuzaixiaoyuan/pkg/database"
	"wobuzaixiaoyuan/pkg/logServer"
	"wobuzaixiaoyuan/pkg/wzxy"
)

// TEST GITHUB ACTION
func main() {
	//setTime()
	log.SetFlags(log.Ltime | log.Ldate)
	logServer.StartLogServer()
	users, err := database.GetUsers()
	eventMap := make(map[string]int)
	if err != nil {
		fmt.Println(err)
		return
	}
	c := cron.New()
	updateSpec := "0 30 * * * *"
	checkinSpec := "0 * * * * *"
	resetSpec := "0 0 0 * * *"
	c.AddFunc(updateSpec, func() {
		log.Printf("更新用户数据")
		dataTmp, err := database.GetUsers()
		if err != nil {
			fmt.Println(err)
		} else {
			users = dataTmp
		}
	})
	c.AddFunc(checkinSpec, func() {
		var checkinUsers []*wzxy.User
		results := make(chan string, 10)
		for _, user := range users {
			// 如果当前时间在用户设置的时间范围内，并且用户今天还没有打过卡
			if !common.CompareTime(user.Start) && common.CompareTime(user.End) && eventMap[user.RealName] < 2 {
				checkinUsers = append(checkinUsers, &wzxy.User{
					RealName: user.RealName,
					Username: user.Username,
					Password: user.Password,
					Result:   results,
				})
				eventMap[user.RealName]++
			}
		}
		if len(checkinUsers) == 0 {
			log.SetOutput(os.Stdout)
			log.Printf("没有需要打卡的用户")
		}
		if len(checkinUsers) != 0 {
			doWork(checkinUsers)
		}
	})
	c.AddFunc(resetSpec, func() {
		log.Printf("重置打卡次数")
		for _, user := range users {
			eventMap[user.RealName] = 0
		}
	})
	c.Start()
	select {}
}
func doWork(users []*wzxy.User) {
	var logFile bytes.Buffer
	multiWriter := io.MultiWriter(os.Stdout, &logFile)
	log.SetOutput(multiWriter)
	var wg sync.WaitGroup
	start := time.Now()
	maxConcurrent := 10
	// 用来控制最大并发数量
	sem := make(chan struct{}, maxConcurrent)
	// 用来接收已签到的wzxy对象
	successCh := make(chan *wzxy.Session, len(users))
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
	data := logFile.String()
	database.UploadLog(data)
	logFile.Reset()
}

func setTime() {
	var cstZone = time.FixedZone("CST", 8*3600) // 东八
	time.Local = cstZone
}
