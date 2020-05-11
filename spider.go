package main

import (
	"./sqlgo"
	"database/sql"
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"os"
	"strings"
)

var conn *sql.DB

func main() {
	conn = sqlgo.InitMySql()

	sqlCmd := fmt.Sprintf(`
								create table IF NOT EXISTS %s(
								bid int auto_increment,
								name varchar(20),
								href varchar(100),
								primary key(bid))
								engine=InnoDB default charset=utf8
								`,"总目录")
	sqlgo.CurdSql(conn, sqlCmd)


	url := "https://www.jd.com/allSort.aspx"
	//url := "https://movie.douban.com/tag/#/?sort=T&tags=%E8%B6%85%E7%BA%A7%E8%8B%B1%E9%9B%84"
	visitURL(url)
	fmt.Println("the end ")
	conn.Close()
}

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
	parseForm(n)
	//do dfs
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		getURL(c)
	}
}

func parseForm(n *html.Node) {
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


		sqlCmd := fmt.Sprintf(`
								create table IF NOT EXISTS %s(
								mid int auto_increment,
								name varchar(20),
								href varchar(100),
								primary key(mid))
								engine=InnoDB default charset=utf8
								`, bigName)

		//sqlCmd := fmt.Sprintf(`drop table %s`, bigName)
		sqlgo.CurdSql(conn, sqlCmd)




		sqlCmd = fmt.Sprintf(`insert into 总目录
									(name)
									values
									('%s')
									`, bigName)
		sqlgo.CurdSql(conn, sqlCmd)


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

			sqlCmd := fmt.Sprintf(`
								create table IF NOT EXISTS %s(
								sid int auto_increment,
								name varchar(20),
								href varchar(100),
								primary key(sid))
								engine=InnoDB default charset=utf8
								`, midName)


			//sqlCmd := fmt.Sprintf(`drop table %s`, bigName)
			sqlgo.CurdSql(conn, sqlCmd)

			sqlCmd = fmt.Sprintf(`insert into %s
									(name,href)
									values
									('%s','%s')
									`, bigName,midName,midHref)

			//sqlCmd := fmt.Sprintf(`drop table %s`, nodeList[0].FirstChild.Data)
			sqlgo.CurdSql(conn, sqlCmd)


			ddNode := dtNode.Parent.Parent
			for q := ddNode.FirstChild; q != nil; q = q.NextSibling {
				if q.Data == "dd" {
					for c := q.FirstChild; c != nil; c = c.NextSibling {
						if c.Data == "a" {
							smlName = c.FirstChild.Data
							smlName = strings.Replace(smlName," ","",-1)
							smlName = strings.Replace(smlName,"/","",-1)
							smlHref = c.Attr[0].Val

							sqlCmd := fmt.Sprintf(`insert into %s
									(name,href)
									values
									('%s','%s')
									`, midName,smlName,smlHref)

							//sqlCmd := fmt.Sprintf(`drop table %s`, nodeList[0].FirstChild.Data)
							sqlgo.CurdSql(conn, sqlCmd)



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
