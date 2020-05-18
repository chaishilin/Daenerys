package main

import (
	"./redis"
	"./sqlgo"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
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
	server := &http.Server{
		Addr:    ":18080",
		Handler: mux,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}


}

func classHandler(w http.ResponseWriter, r *http.Request) {

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

}

func logHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		file,_ := ioutil.ReadFile("./root/test.html")
		fmt.Fprintf(w,string(file))
	}else if r.Method == "POST" {
		r.ParseForm()
		conn := redisconfirm.InitRedis(passwd)
		logCheck := redisconfirm.LogCheck(&conn, r.Form.Get("name"), r.Form.Get("pwd"))
		fmt.Fprint(w,logCheck)

	}

	/*
	if r.Method == "POST" {
		r.ParseForm()
		conn := redisconfirm.InitRedis(passwd)
		logCheck := redisconfirm.LogCheck(&conn, r.Form.Get("username"), r.Form.Get("passwd"))

		if logCheck == false {
			t, _ := template.ParseFiles("./root/log.html")
			t.Execute(w, "请输入正确的用户名和密码")
		} else {
			t, _ := template.ParseFiles("./root/hello.html")
			t.Execute(w, 0)
		}

	} else if r.Method == "GET" {
		t, _ := template.ParseFiles("./root/log.html")
		t.Execute(w, "")
	}

	 */
}

func registHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		conn := redisconfirm.InitRedis(passwd)
		r.ParseForm()
		_, state, _ := redisconfirm.Register(&conn, r.Form.Get("username"), r.Form.Get("passwd"))
		if state == redisconfirm.SetOk {
			t, _ := template.ParseFiles("./root/log.html")
			t.Execute(w, "")
		} else {
			t, _ := template.ParseFiles("./root/regist.html")
			t.Execute(w, "用户名已被占用")
		}
	} else if r.Method == "GET" {
		t, _ := template.ParseFiles("./root/regist.html")
		t.Execute(w, "")

	}
}
