package webTemplate

func GetLogTmpl() string {
	return `<html>
    <head>
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <title>{{.Time}}</title>
        <style>
        body {
                font-size: 18px;
                text-align: center;
                margin: 0;
                padding: 0;
            }
            h3 {
                font-size: 24px;
                margin-top: 20px;
                margin-bottom: 20px;
            }
            textarea {
                margin: 10px auto;
                padding: 10px;
                border: 1px solid black;
                border-radius: 5px;
                box-sizing: border-box;
                width: 90%;
                max-width: 600px;
                height: 300px;
                resize: none;
                font-size: 16px;
                line-height: 1.5;
                font-family: 'Arial', sans-serif;
            }
            </style>
    </head>
    <body>
        <h3>{{.Time}}</h3>
        <textarea rows='10' style="width: 100%;">{{.Content}}</textarea>
    </body>
</html>`
}
