package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type classes struct {
	Name string
	Id int
	Pid int
	Href string
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./root/hello.html")
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprint(w, err)
		return
	}
	classlist := []classes{
		{Name: "数码产品",Id: 1,Pid: 0,Href: "www.shuma.com"},
		{Name: "手机",Id: 2,Pid: 1,Href: "https://channel.jd.com/shouji.html"},
		{Name: "苹果手机",Id: 3,Pid: 1,Href: "https://list.jd.com/list.html?cat=652,829,845"},
		{Name: "相机",Id: 4,Pid: 1,Href: "www.camera.com"},
		{Name: "索尼相机",Id: 5,Pid: 4,Href: "www.sonycamera.com"},
	}
	t.Execute(w, classlist)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", indexHandler)

	server := &http.Server {
		Addr:    ":8080",
		Handler: mux,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
