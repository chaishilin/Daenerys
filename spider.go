package main

import (
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"os"
)

func main() {
	url := "https://www.jd.com/allSort.aspx"
	//url := "https://movie.douban.com/tag/#/?sort=T&tags=%E8%B6%85%E7%BA%A7%E8%8B%B1%E9%9B%84"
	visitURL(url)
	fmt.Println("the end ")
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
	result := haveLabel(n, []string{"category-item m", "mt", "item-title", "span"}, &nodeList)
	if result == true {
		fmt.Println("\n----------------------------------------------\n")
		fmt.Printf("大类: %s\t\n", nodeList[0].FirstChild.Data)
	}

	nodeList = []*html.Node{}
	result = haveLabel(n, []string{"category-item m", "mc", "items", "dl", "dt", "a"}, &nodeList)

	if result == true {
		for _, dtNode := range nodeList {
			fmt.Printf("--中类: %s\t", dtNode.FirstChild.Data)
			fmt.Printf("链接: %s \n", dtNode.Attr[0].Val)
			ddNode := dtNode.Parent.Parent
			for q := ddNode.FirstChild; q != nil; q = q.NextSibling {
				if q.Data == "dd" {
					for c := q.FirstChild; c != nil; c = c.NextSibling {
						if c.Data == "a" {
							fmt.Printf("----小类: %s\t", c.FirstChild.Data)
							fmt.Printf("链接: %s \n", c.Attr[0].Val)
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