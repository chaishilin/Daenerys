package main

import (
	"./jdSpider"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {

	url := "https://www.jd.com/allSort.aspx"
	resp, _ := http.Get(url)
	b, _ := ioutil.ReadAll(resp.Body)
	s := fmt.Sprintf("%s", b)
	resp.Body.Close()
	//conn := sqlgo.InitMySql()
	//sqlgo.CreateTables(conn)
	jdSpider.DoSpider(s,nil)


}
