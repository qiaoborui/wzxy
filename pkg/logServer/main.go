package logServer

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"wobuzaixiaoyuan/pkg/database"
)

// @Title startLogServer
// @Description 启动日志服务器
func StartLogServer() {
	shanghaiLoc := time.FixedZone("UTC+8", 8*60*60)

	go func() {
		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}

		host := os.Getenv("HOST")
		if host == "" {
			host = "0.0.0.0"
		}

		serverAddress := fmt.Sprintf("%s:%s", host, port)
		fmt.Printf("Server listening on http://%s\n", serverAddress)

		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "<html><head><title>日志查看</title></head><body><h1>Log Viewer</h1>")
			logs, err := database.GetLogs()
			if err != nil {
				fmt.Fprintf(w, "Failed to retrieve logs: %s<br>", err)
				http.Error(w, "Failed to retrieve logs", http.StatusInternalServerError)
				return
			}
			numLogs := len(logs)
			fmt.Fprintf(w, "共有: %d个日志<br>", numLogs)
			for _, log := range logs {
				createTime := log.CreatedAt.In(shanghaiLoc).Format("2006-01-02 15:04:05")
				logID := log.Object.ID
				fmt.Fprintf(w, "<a href='/logs/%s'>%s</a><br>", logID, createTime)
			}
			fmt.Fprintf(w, "</body></html>")
		})

		http.HandleFunc("/logs/", func(w http.ResponseWriter, r *http.Request) {
			logID := r.URL.Path[len("/logs/"):]
			logContent, err := database.GetSpecficLog(logID)
			if err != nil {
				http.Error(w, "File "+logID+" not found.", http.StatusNotFound)
				return
			}
			createTime := logContent.CreatedAt.In(shanghaiLoc).Format("2006-01-02 15:04:05")
			logContentText := logContent.Content
			fmt.Fprintf(w, "<html><head><meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\"><title>%s</title></head><body><h3>%s</h3>", createTime, createTime)
			fmt.Fprintf(w, "<textarea rows='10' style=\"width: 100%%;\">%s</textarea>", logContentText)
			fmt.Fprintf(w, "</body></html>")
		})

		err := http.ListenAndServe(serverAddress, nil)
		if err != nil {
			log.Fatal(err)
		}
	}()
}
