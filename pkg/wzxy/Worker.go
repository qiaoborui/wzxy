package wzxy

import (
	"bytes"
	"io"
	"log"
	"os"
	"sync"
	"time"
	"wobuzaixiaoyuan/pkg/database"
)

func DoWork(users []*User) {
	var logFile bytes.Buffer
	multiWriter := io.MultiWriter(os.Stdout, &logFile)
	log.SetOutput(multiWriter)
	var wg sync.WaitGroup
	start := time.Now()
	maxConcurrent := 10
	// 用来控制最大并发数量
	sem := make(chan struct{}, maxConcurrent)
	// 用来接收已签到的wzxy对象
	successCh := make(chan *Session, len(users))
	failedCh := make(chan *Session, len(users))
	for _, user := range users {
		w := NewSession(user)
		sem <- struct{}{}
		wg.Add(1)
		go func(w *Session) {
			defer func() { <-sem }()
			defer wg.Done()
			w.Login()
			tasks := w.GetSignList()
			w.Sign(tasks)
			if w.Err != nil {
				failedCh <- w
				return
			}
			successCh <- w
		}(w)
	}
	wg.Wait()
	close(sem)
	close(successCh)
	close(failedCh)
	// 按顺序输出已签到的wzxy对象
	for w := range successCh {
		log.Printf("%s\n", <-w.User.Result)
	}
	for w := range failedCh {
		log.Printf("[%s] occur errors: %v",w.User.RealName,w.Err)
	}
	elapsed := time.Since(start)
	log.Printf("程序运行时间为：%s \n", elapsed)
	data := logFile.String()
	err := database.UploadLog(data)
	if err != nil {
		log.Printf("上传日志失败: %v\n", err)
	}
	logFile.Reset()
}
