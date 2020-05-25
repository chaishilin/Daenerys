package main

import (
	"./jdSpider"
	"./sqlgo"
	"fmt"
	"io/ioutil"
	"net/http"
)

type classInfo struct {
	Class_id int
	Class_name string
	class_href string
}


func main() {

	url := "https://www.jd.com/allSort.aspx"
	resp, _ := http.Get(url)
	b, _ := ioutil.ReadAll(resp.Body)
	s := fmt.Sprintf("%s", b)
	resp.Body.Close()



	sqlgo.DelAll()
	sqlgo.CreateTables()
	jdSpider.DoSpider(s)

	jdSpider.GetGoodsIntr()

	/*
	result := sqlgo.SelectSql(conn,"select * from classTable")
	for _,v := range result{
		fmt.Printf("id : %s ,名字 :  %s,链接 : %s \n",v[0],v[1],v[2])
	}
	*/

}
