package jdSpider

import (
	"../sqlgo"
	"database/sql"
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"os"
	"strings"
	"sync"
)

var conn *sql.DB
//var count = 1
var wg sync.WaitGroup
var mu sync.Mutex

/*
func main() {
	conn = sqlgo.InitMySql()
	sqlgo.CreateTables(conn)

	url := "https://www.jd.com/allSort.aspx"
	visitURL(url)
	wg.Wait()
	fmt.Println("the end ")
	conn.Close()
}

 */

func visitURL(url string) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch get %v", err)
		os.Exit(1)
	}
	doc, err := html.Parse(resp.Body)
	resp.Body.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch parse%v", err)
		os.Exit(1)
	}
	getURL(doc)
}

func getURL(n *html.Node) {
	wg.Add(1)
	go parseForm(n)
	//do dfs
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		getURL(c)
	}
}

func parseForm(n *html.Node) {
	defer wg.Done()
	nodeList := []*html.Node{}
	var bigName string
	var midName string
	var midHref string
	var smlName string
	var smlHref string

	result := haveLabel(n, []string{"category-item m", "mt", "item-title", "span"}, &nodeList)
	if result == true {
		bigName = nodeList[0].FirstChild.Data
		fmt.Println("\n----------------------------------------------\n")
		fmt.Printf("大类: %s\t\n", bigName)
		mu.Lock()
		count++
		sqlCmd1 := fmt.Sprintf(`insert into classTable
									(class_id,class_name)
									values
									(%d,'%s')
									`, count,bigName)
		sqlCmd2 := fmt.Sprintf(`insert into classRelate
									(class_id,pid,class_name)
									values
									(%d,0,'%s')
									`, count,bigName)
		mu.Unlock()


		sqlgo.CurdSql(conn, sqlCmd1)

		sqlgo.CurdSql(conn, sqlCmd2)


	}

	nodeList = []*html.Node{}
	result = haveLabel(n, []string{"category-item m", "mc", "items", "dl", "dt", "a"}, &nodeList)

	if result == true {
		for _, dtNode := range nodeList {
			midName = dtNode.FirstChild.Data
			midName = strings.Replace(midName," ","",-1)
			midName = strings.Replace(midName,"/","",-1)
			midHref = dtNode.Attr[0].Val
			fmt.Printf("--中类: %s\t", midName)
			fmt.Printf("链接: %s \n", midHref)
			mu.Lock()
			count++
			sqlCmd1 := fmt.Sprintf(`insert into classTable
									(class_name,class_href,class_id)
									values
									('%s','%s',%d)
									`, midName,midHref,count)

			sqlCmd2 := fmt.Sprintf(`insert into classRelate
									(class_id,pid,class_name)
									values
									(%d,0,'%s')
									`, count,midName)
			mu.Unlock()


			sqlgo.CurdSql(conn, sqlCmd1)



			sqlgo.CurdSql(conn, sqlCmd2)


			ddNode := dtNode.Parent.Parent
			for q := ddNode.FirstChild; q != nil; q = q.NextSibling {
				if q.Data == "dd" {
					for c := q.FirstChild; c != nil; c = c.NextSibling {
						if c.Data == "a" {
							smlName = c.FirstChild.Data
							smlName = strings.Replace(smlName," ","",-1)
							smlName = strings.Replace(smlName,"/","",-1)
							smlHref = c.Attr[0].Val
							mu.Lock()
							count++
							sqlCmd1 := fmt.Sprintf(`insert into classTable
									(class_name,class_href,class_id)
									values
									('%s','%s',%d)
									`, smlName,smlHref,count)
							sqlCmd2 := fmt.Sprintf(`insert into classRelate
									(class_id,pid,class_name)
									values
									(%d,0,'%s')
									`, count,smlName)
							mu.Unlock()


							sqlgo.CurdSql(conn, sqlCmd1)



							sqlgo.CurdSql(conn, sqlCmd2)



							fmt.Printf("----小类: %s\t", smlName)
							fmt.Printf("链接: %s \n", smlHref)

						}
					}
				}
			}
		}
	}
}

func haveLabel(n *html.Node, labelList []string, nodeList *[]*html.Node) bool {
	if len(labelList) == 1 {
		if n.Data == labelList[0] {
			*nodeList = append(*nodeList, n)
			return true
		}
		return false
	}
	flag := 0
	for _, v := range n.Attr {
		if v.Key == "class" && v.Val == labelList[0] {
			flag = 1
		}
	}
	if n.Data == labelList[0] {
		flag = 1
	}
	if flag == 1 {
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			haveLabel(c, labelList[1:], nodeList)
		}
	}
	if len(*nodeList) > 0 {
		return true
	} else {
		return false
	}
}
