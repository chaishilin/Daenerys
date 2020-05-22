package main

import (
	"io"

	"os"
)


func main(){
	/*
	href :="https://e.jd.com/ebook.html"
	resp,err:=http.Get(href)
	fmt.Println(err)
	bBody,_ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(bBody))

	 */
	/*
	分类和商品信息
	select b.class_name,a.class_id,a.goods_id,a.goods_name
	from goodsTable as a join classTable as b
	where a.class_id = b.class_id and a.class_id = 212;
	 */

	s := []string{"123421\n","asfdfds\n","2dw"}
	file,_ := os.OpenFile("./out.txt",os.O_WRONLY|os.O_CREATE|os.O_APPEND,0666)
	for _,ss:= range s{
		io.WriteString(file,ss)

	}
	/*
	var pwg sync.WaitGroup
	//href := "https://list.jd.com/list.html?cat=6144,12041,12049"
	href := "https://list.jd.com/list.html?cat=1672,2615,9186"
	pwg.Add(1)
	conn := sqlgo.InitMySql()

	sqlgo.DelAll(conn)

	sqlgo.CreateTables(conn)
	go jdSpider.GetGood(href,0,&pwg)
	pwg.Wait()
	fmt.Println("over")

	 */



}