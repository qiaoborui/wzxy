package webTemplate

func GetIndexTemplate() string {
	return `<!DOCTYPE html>
<html>
    <head>
        <title>日志查看</title>
        <style>
             body {
                    font-size: 18px;
                    text-align: center;
                }
                h1 {
                    font-size: 32px;
                    margin-top: 20px;
                    margin-bottom: 20px;
                }
                .col-12 {
                    margin: 10px auto;
                    padding: 10px;
                    border: 1px solid black;
                    border-radius: 5px;
                    text-align: center;
                }
                a {
                    text-decoration: none;
                    color: blue;
                }
        </style>
    </head>
    <body>
        <h1>日志查看</h1>
        共有: {{.Num}}个日志<br>
        {{ range $log := .Logs}}
              <div class="col-12">
                  <a href="/logs/{{$log.LogID}}"> {{$log.CreateTime}} </a>
              </div>
        {{end}}
    </body>
</html>`
}
