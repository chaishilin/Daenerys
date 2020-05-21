package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main(){
	href :="//item.jd.com/3650111.html"
	resp,_:=http.Get(href)
	bBody,_ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(bBody))
}