package jdSpider

import (
	"../sqlgo"
	"fmt"
	"regexp"
	"sync"
)

const (
	itemReg  = `<!--  /widget/cat-item/cat-item.vm -->[\s\S]+?<!--/ /widget/cat-item/cat-item.vm -->`
	titleReg = `<span>(.+)</span>`
	dlReg    = `<dl class=\"clearfix\"[\s\S]*?</dl>`
	dtReg    = `<dt>[\s\S]+?</dt>`
	ddReg    = `<dd>[\s\S]+?</dd>`
	infoReg  = `<a href="([\/\.\?\-,=&_0-9a-zA-Z]+)"[^>]*>([^<]+)</a>`
)

var count = 0

func DoSpider(htmlMsg string) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	items := regexp.MustCompile(itemReg).FindAllString(htmlMsg, -1)

	for _, each := range items {
		title := regexp.MustCompile(titleReg)
		titles := title.FindAllSubmatch([]byte(each), -1)

		count++
		fmt.Println("--", count, 0, string(titles[0][1]))
		sqlgo.InsertClass(count,string(titles[0][1]),"")
		sqlgo.InsertRelate(count,0)

		pid := count
		dls := regexp.MustCompile(dlReg).FindAllString(each, -1)

		for _, eachdl := range dls {
			func(eachdl string) {
				dtresults, pid := findMatch(eachdl, dtReg, pid)
				for _, v := range dtresults {
					count++
					fmt.Println("----", count, pid, v[1], v[0])
					sqlgo.InsertClass(count,v[1],v[0])
					sqlgo.InsertRelate(count,pid)
				}
				dlresults, pid := findMatch(eachdl, ddReg, count)
				for _, v := range dlresults {
					count++
					href := v[0]
					if href[:6] != "https:"{
						href =  "https:"+href
					}
					//fmt.Println("--------", count, pid, v[1], href)
					sqlgo.InsertClass(count,v[1],href)
					sqlgo.InsertRelate(count,pid)
					wg.Add(1)
					go GetGood(href,count,&wg,&mu)
				}
			}(eachdl)
		}
	}
	wg.Wait()
	fmt.Println(count)
}

func findMatch(str string, reg string, pid int) ([][]string, int) {
	var result [][]string
	regex := regexp.MustCompile(reg)
	regexs := regex.FindAllString(str, -1)
	for _, eachItem := range regexs {
		info := regexp.MustCompile(infoReg)
		infos := info.FindAllSubmatch([]byte(eachItem), -1)
		for _, eachResult := range infos {
			result = append(result, []string{string(eachResult[1]), string(eachResult[2])})
		}
	}
	return result, pid
}


