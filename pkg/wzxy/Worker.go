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

type Result struct {
	RealName string
	Status   int
	Msg      string
}

func DoWork(users []*User) []Result {
	start := time.Now()
	_, logFile := setupLogging()
	maxConcurrent := 10
	sem := make(chan struct{}, maxConcurrent)
	successCh := make(chan *Session, len(users))
	createWorkers(users, sem, successCh)
	results := make([]Result, 0)
	// 按顺序输出已签到的wzxy对象
	for w := range successCh {
		result := <-w.User.Result
		res := Result{}
		res.RealName = w.User.RealName
		res.Msg = result
		if bytes.Contains([]byte(result), []byte("成功")) {
			res.Status = 0
		} else {
			res.Status = 1
		}
		results = append(results, res)
		log.Printf("%s\n", result)
	}
	elapsed := time.Since(start)
	log.Printf("程序运行时间为：%s \n", elapsed)
	data := logFile.String()
	err := database.UploadLog(data)
	if err != nil {
		log.Printf("上传日志失败: %v\n", err)
	}
	logFile.Reset()
	return results
}

func setupLogging() (io.Writer, *bytes.Buffer) {
	var logFile bytes.Buffer
	multiWriter := io.MultiWriter(os.Stdout, &logFile)
	log.SetOutput(multiWriter)
	return multiWriter, &logFile
}

func createWorkers(users []*User, sem chan struct{}, successCh chan *Session) {
	var wg sync.WaitGroup
	for _, user := range users {
		w := NewSession(user)
		sem <- struct{}{}
		wg.Add(1)
		go func(w *Session) {
			defer func() { <-sem }()
			defer wg.Done()
			if err := w.Login(); err != nil {
				log.Printf("[%s] login failed: %v\n", w.User.RealName, err)
				return
			}
			_ = w.Sign()
			successCh <- w
		}(w)
	}
	wg.Wait()
	close(sem)
	close(successCh)
}
