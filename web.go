package main

import (
	"./redis"
	"./sqlgo"
	"fmt"
	"github.com/dchest/captcha"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type classes struct {
	Name string
	Id   int
	Pid  int
	Href string
}

const (
	queryById = `select a.class_id,a.class_name,a.class_href 
					from classTable as a join classRelate as b 
					where a.class_id = b.class_id and 
					(b.pid = ? or b.class_id = ?);`
	queryByName = `select a.class_id,a.class_name,a.class_href 
						from classTable as a join classRelate as b 
						where a.class_id = b.class_id and 
						a.class_name regexp '%s';`
)

var passwd string

func main() {
	passwd = os.Args[1:][0]
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", classHandler)
	mux.HandleFunc("/", logHandler)
	mux.HandleFunc("/regist", registHandler)
	mux.HandleFunc("/process", captchaVerify)
	mux.HandleFunc("/captcha/newId", newCapId)
	mux.HandleFunc("/template", makeTemplate)
	mux.Handle("/captcha/", captcha.Server(captcha.StdWidth, captcha.StdHeight))
	server := &http.Server{
		Addr:    ":18080",
		Handler: mux,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func newCapId(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, captcha.New())
}

func captchaVerify(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html;charset=utf-8")
	if !captcha.VerifyString(r.FormValue("captchaId"), r.FormValue("captchaSolution")) {
		//每个CaptchaId只能校验一次
		fmt.Fprint(w, "false")
	} else {
		fmt.Fprint(w, "true")
	}
}

func makeTemplate(w http.ResponseWriter, r *http.Request)  {
	file,_ := ioutil.ReadFile("./root/classTemplate.html")
	fmt.Fprint(w,string(file))

}

func classHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		file, _ := ioutil.ReadFile("./root/test2.html")
		fmt.Fprintf(w, string(file))
	} else if r.Method == "POST" {
		r.ParseForm()
		fmt.Fprint(w, 3)

	}
	/*
		r.ParseForm()

		classList := []classes{}
		conn := sqlgo.InitMySql()
		//findAll := "select * from classTable"
		query := r.Form.Get("input")
		_, err := strconv.Atoi(query)
		var result [][]string
		if err != nil {
			if len(query) > 0 {
				result = sqlgo.SelectSql(conn, fmt.Sprintf(queryByName, fmt.Sprintf(".*%s.*", query))) //"'."+query+".'"
			}
		} else {
			result = sqlgo.SelectSql(conn, queryById, query, query)
		}
		for _, v := range result {
			intId, _ := strconv.Atoi(v[0])
			classList = append(classList, classes{Id: intId, Name: v[1], Href: v[2]})
		}

		t, err := template.ParseFiles("./root/hello.html")
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprint(w, err)
			return
		}

		t.Execute(w, classList)

	*/

}

func logHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		file, _ := ioutil.ReadFile("./root/test.html")
		fmt.Fprintf(w, string(file))
	} else if r.Method == "POST" {
		r.ParseForm()
		conn := redisconfirm.InitRedis(passwd)
		logCheck := redisconfirm.LogCheck(&conn, r.Form.Get("name"), r.Form.Get("pwd"))
		fmt.Fprint(w, logCheck)
	}
}

func registHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		conn := redisconfirm.InitRedis(passwd)
		r.ParseForm()
		username := r.Form.Get("username")
		passwd := r.Form.Get("passwd")
		email_addr := r.Form.Get("email")
		_, state, _ := redisconfirm.Register(&conn, username, passwd)
		if state == redisconfirm.SetOk {
			fmt.Fprintf(w, "ok")
			//使用mysql存入用户邮箱、用户名，密码
			conn := sqlgo.InitMySql()
			sqlgo.InsertUser(conn,username,email_addr)
		} else {
			fmt.Fprintf(w, "exist")
		}
	} else if r.Method == "GET" {
		t, _ := template.ParseFiles("./root/regist.html")
		t.Execute(w, "")

	}
}

func doVerify() {

}
