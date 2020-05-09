package main

import (
	"fmt"
	"golang.org/x/net/html"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	end := make(chan int)
	url := "https://www.jd.com/allSort.aspx"
	urllist := []string{url}
	count := 0
	for _, v := range urllist {
		go visitURL(v, end)
		count++
	}
	if count > 0 {
		for range end {
			count--
			if count == 0 {
				break
			}
		}
	}
	fmt.Println("the end ")
}

func visitURL(url string, end chan int) {
	defer func() { end <- 1 }()
	resp, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch :%v\n", err)
		os.Exit(1)
	}
	b, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch : reading %v\n", err)
		os.Exit(1)
	}

	node, err := html.Parse(resp.Body)

	resp.Body.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch : reading %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("%v\n", node)
	fmt.Printf("get : %s  %s\n", url, b[:10])

}
