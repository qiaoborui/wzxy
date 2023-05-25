package main

import (
	"fmt"
	"sync"
	"time"
	"wobuzaixiaoyuan/utils"
	"wobuzaixiaoyuan/wzxy"
)

// TEST GITHUB ACTION
func main() {
	start := time.Now()
	var wg sync.WaitGroup
	storage, err := utils.FetchData("{\"status\": 1 }")
	if err != nil {
		fmt.Println(err)
		return
	}
	results := make(chan string, 10)
	var users []*wzxy.User
	if len(storage.Results) == 0 {
		fmt.Println("没有需要打卡的用户")
		return
	}
	for _, user := range storage.Results {
		users = append(users, &wzxy.User{
			RealName: user.RealName,
			Username: user.Username,
			Password: user.Password,
			Result:   results,
		})
	}
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
		fmt.Printf("%s\n", <-w.User.Result)
	}
	elapsed := time.Since(start)
	fmt.Printf("程序运行时间为：%s", elapsed)
}
