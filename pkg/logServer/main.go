package logServer

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
	"wobuzaixiaoyuan/pkg/database"
	"wobuzaixiaoyuan/pkg/webTemplate"
)

type ATag struct {
	CreateTime string
	LogID      string
}

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
			logs, err := database.GetLogs()
			if err != nil {
				http.Error(w, "获取日志失败", http.StatusInternalServerError)
				return
			}
			numLogs := len(logs)
			tmpl, err := template.New("index").Parse(webTemplate.GetIndexTemplate())
			if err != nil {
				http.Error(w, "解析模板失败", http.StatusInternalServerError)
				return
			}
			var aTags []ATag
			for i := 0; i < numLogs; i++ {
				createTime := logs[i].CreatedAt.In(shanghaiLoc).Format("2006-01-02 15:04:05")
				aTags = append(aTags, ATag{
					CreateTime: createTime,
					LogID:      logs[i].Object.ID,
				})
			}
			err = tmpl.Execute(w, map[string]interface{}{
				"Num":  numLogs,
				"Logs": aTags,
			})
			if err != nil {
				http.Error(w, fmt.Sprintf("渲染模板失败 err: %v", err), http.StatusInternalServerError)
				return
			}
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
			tmpl, err := template.New("log").Parse(webTemplate.GetLogTmpl())
			if err != nil {
				http.Error(w, "解析模板失败", http.StatusInternalServerError)
				return
			}
			err = tmpl.Execute(w, map[string]interface{}{
				"Time":    createTime,
				"Content": logContentText,
			})
			if err != nil {
				http.Error(w, fmt.Sprintf("渲染模板失败 err: %v", err), http.StatusInternalServerError)
				return
			}
		})

		err := http.ListenAndServe(serverAddress, nil)
		if err != nil {
			log.Fatal(err)
		}
	}()
}
