package main

import (
	"encoding/json"
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const logFile = "logs.json"

//Log　掲示板に保存するデータを構造体で定義
type Log struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Body  string `json:"body"`
	CTime int64  `json:"ctime"`
}

func main() {
	println("server - http://localhost:8080")
	//URIに対応するハンドラを登録
	http.HandleFunc("/", showHandler)
	http.HandleFunc("/write", writeHandler)
	//サーバーを起動
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("error starting http server:", err)
		return
	}
}

//書き込みログを画面に表示する
func showHandler(w http.ResponseWriter, r *http.Request) {
	//ログを読み出してHTMLを生成
	htmlLog := ""
	logs := loadLogs()
	for _, i := range logs {
		htmlLog += fmt.Sprintf(
			"<p>(%d) <span>%s</span>: %s --- %s</p>",
			i.ID,
			html.EscapeString(i.Name),
			html.EscapeString(i.Body),
			time.Unix(i.CTime, 0).Format("2006/1/2 15:04"),
		)
	}
	//HTML全体を出力
	htmlBody := "<html><head><style>" +
		"p{border: 1px solid silver; padding: 1em;}" +
		"span {background-color: #eef;}" +
		"</style></head><body><h1>BBS</h1>" +
		getForm() + htmlLog + "</body></html>"
	w.Write([]byte(htmlBody))
}

//フォームから送信された内要を書き込み
func writeHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() //フォームを解析
	var log Log
	log.Name = r.Form["name"][0]
	log.Body = r.Form["body"][0]
	if log.Name == "" {
		log.Name = "名無し"
	}
	logs := loadLogs() //既存のデータを読み出し
	log.ID = len(logs) + 1
	log.CTime = time.Now().Unix()
	logs = append(logs, log)      //追記
	saveLogs(logs)                //保存
	http.Redirect(w, r, "/", 302) //リダイレクト
}

//書き込みフォームを返す
func getForm() string {
	return "<div><form action='/write' method='POST'>" +
		"名前: <input type='text' name='name'><br>" +
		"本文: <input type='text' name='body' style='width:30em;'><br>" +
		"<input type='submit' value='書込'>" +
		"</form></div><hr>"
}

//ファイルからログファイルの読み込み
func loadLogs() []Log {
	//ファイルを開く
	text, err := ioutil.ReadFile(logFile)
	if err != nil {
		return make([]Log, 0)
	}
	//JSONをパース
	var logs []Log
	json.Unmarshal([]byte(text), &logs)
	return logs
}

//ログファイルの書き込み
func saveLogs(logs []Log) {
	//JSONにエンコード
	bytes, _ := json.Marshal(logs)
	//ファイルに書き込む
	ioutil.WriteFile(logFile, bytes, 0644)
}
