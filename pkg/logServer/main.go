package logServer

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

// @Title startLogServer
// @Description 启动日志服务器
func StartLogServer() {
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		err = os.Mkdir("logs", 0755)
		if err != nil {
			log.Fatal(err)
		}
	}
	go func() {
		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}

		host := os.Getenv("wzxy_HOST")
		if host == "" {
			host = "0.0.0.0"
		}
		fmt.Printf("Server listening on http://%s:%s\n", host, port)
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "<html><head><title>我不在校园</title></head><body><h1>日志查看</h1>")
			logDir := "./logs"
			files, err := os.Open(logDir)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer files.Close()

			logFiles, err := files.Readdir(-1)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			fmt.Fprintf(w, "共有 %d 个日志<br>", len(logFiles))
			for _, file := range logFiles {
				if file.Mode().IsRegular() {
					//filePath := filepath.Join(logDir, file.Name())
					fmt.Fprintf(w, "<a href='/logs/%s'>%s</a><br>", file.Name(), file.Name())
				}
			}

			fmt.Fprintf(w, "</body></html>")
		})

		http.HandleFunc("/logs/", func(w http.ResponseWriter, r *http.Request) {
			logFile := r.URL.Path[len("/logs/"):]
			logFilePath := fmt.Sprintf("logs/%s", logFile)
			contents, err := os.ReadFile(logFilePath)
			if err != nil {
				http.Error(w, "File "+logFile+" not found.", 404)
				return
			}
			fmt.Fprintf(w, "<html><head><meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\"><title>%s</title></head><body><h3>%s</h3>", logFile, logFile)
			fmt.Fprintf(w, "<textarea rows='10' style=\"width: 100%%;\">%s</textarea>", contents)
			fmt.Fprintf(w, "</body></html>")
		})

		http.ListenAndServe(":8080", nil)
	}()
}
