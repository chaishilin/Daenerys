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

func classHandler(w http.ResponseWriter, r *http.Request) {

	classList := []classes{}
	r.ParseForm()

	conn := sqlgo.InitMySql()
	//findAll := "select * from classTable"

	query := r.Form.Get("input")
	_, err := strconv.Atoi(query)
	var result [][]string
	if err != nil {
		if len(query) > 0{
		result = sqlgo.SelectSql(conn, fmt.Sprintf(queryByName,fmt.Sprintf(".*%s.*",query)))//"'."+query+".'"
		}
	} else {
		result = sqlgo.SelectSql(conn, queryById, query, query)
	}
	for _, v := range result {
		intId, _ := strconv.Atoi(v[0])
		classList = append(classList, classes{Id: intId, Name: v[1], Href: v[2]})
		//fmt.Printf("id : %s ,名字 :  %s,链接 : %s \n",v[0],v[1],v[2])
	}


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

	server := &http.Server{
		Addr:    ":18080",
		Handler: mux,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
