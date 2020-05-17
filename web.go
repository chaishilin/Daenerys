package main

import (
	"./sqlgo"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type classes struct {
	Name string
	Id int
	Pid int
	Href string
}

func classHandler(w http.ResponseWriter, r *http.Request) {

	classList := []classes{}
	r.ParseForm()
	qId := r.Form.Get("input")
	conn := sqlgo.InitMySql()
	//findAll := "select * from classTable"
	findClassId := `select a.class_id,a.class_name,a.class_href 
					from classTable as a join classRelate as b 
					where a.class_id = b.class_id and 
					(b.pid = ? or b.class_id = ?);`
	result := sqlgo.SelectSql(conn,findClassId,qId,qId)
	for _,v := range result{
		intId,_ := strconv.Atoi(v[0])
		classList = append(classList,classes{Id: intId,Name: v[1],Href: v[2]})
		//fmt.Printf("id : %s ,名字 :  %s,链接 : %s \n",v[0],v[1],v[2])
	}
	/*
	如何通过post请求，得到希望查询的商品分类信息？
	 */

	t, err := template.ParseFiles("./root/hello.html")
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprint(w, err)
		return
	}

	t.Execute(w, classList)
}

func main() {


	mux := http.NewServeMux()
	mux.HandleFunc("/", classHandler)

	server := &http.Server {
		Addr:    ":8080",
		Handler: mux,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
