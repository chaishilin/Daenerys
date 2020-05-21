package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"sync"
)

const (
	eachGood = `<li.*class="gl-item">[\s\S]*?</div>[.\s]*</li>`
	goodPrice = `<em>￥</em><i>(.*)</i>`
	goodHref = `<div class="p-name p-name-type-3">[\s\S]*?href="(.*?)"[\s\S]*?</div>`
	goodName = `<div class="p-name p-name-type-3">[\s\S]*<em>(<.*>)?([\S\s]*)</em>[\n.\s\S]*</div>`
	goodPtable = `<div class="Ptable-item">[\s\S]*?</div>`
	goodPackage = `<div class="package-list">[\s\S]*?</div>`
)


var wg sync.WaitGroup
func main() {
	hrefChan := make(chan string)


	href := "https://list.jd.com/list.html?cat=1672,2577,3997"
	wg.Add(1)
	go getClassDetal(href,hrefChan)
	count := 0
	for href = range hrefChan{
		wg.Add(1)
		go getGoodDetal(href)
		count ++
	}
	wg.Wait()
	fmt.Println("end: ",count)


}
func getGoodDetal(href string){
	//href := <- hrefChan
	defer wg.Done()

	if href[:6] != "https:"{
		href =  "https:"+href
	}

	fmt.Println("get",href)
	resp,_:=http.Get(href)
	bBody,_ := ioutil.ReadAll(resp.Body)
	body := string(bBody)
	resp.Body.Close()

	ptableList := regexp.MustCompile(goodPtable).FindAllString(body,-1)
	for _,item := range ptableList{
		paserItem(item)
	}
	packageList := regexp.MustCompile(goodPackage).FindAllString(body,-1)
	if len(packageList) == 0 {
		//fmt.Println("packageList : 0")
	}else{
		paserItem(packageList[0])
	}

}
func getClassDetal(href string,hrefChan chan string){
	defer wg.Done()
	defer close(hrefChan)
	//href := "https://list.jd.com/list.html?cat=9987,830,13661"
	resp,_:=http.Get(href)
	bBody,_ := ioutil.ReadAll(resp.Body)
	body := string(bBody)
	resp.Body.Close()
	//获得分类下的商品页面，得到商品链接、名称、价格
	//还需要进入具体商品，得到商品简介、规格

	goodList := regexp.MustCompile(eachGood).FindAllString(body,-1)
	for _,good:=range goodList{
		fmt.Println("----------------------")
		name:=regexp.MustCompile(goodName)
		nameResult := name.FindAllSubmatch([]byte(good),-1)

		itemName := string(nameResult[0][2])
		itemName = strings.Replace(itemName,"\n","",-1)
		fmt.Println("name: ",itemName)

		price:=regexp.MustCompile(goodPrice)
		priceResult := price.FindAllSubmatch([]byte(good),-1)
		_ =string(priceResult[0][1])
		itemPrice :=string(priceResult[0][1])

		fmt.Println("price: ",itemPrice)

		href:=regexp.MustCompile(goodHref)
		hrefResult := href.FindAllSubmatch([]byte(good),-1)

		itemHref := string(hrefResult[0][1])
		fmt.Println("href: ", itemHref)
		hrefChan<-itemHref

	}
}

func paserItem(str string){
	fmt.Println("=======paserItem==========")
	h3Reg := regexp.MustCompile(`<h3>(.*)</h3>`).FindAllStringSubmatch(str,-1)
	title := h3Reg[0][1]
	fmt.Println("title:",title)
	dlReg := regexp.MustCompile(`<dt>(.*)</dt><dd>(.*)</dd>`).FindAllStringSubmatch(str,-1)
	for _,dl := range(dlReg){
		fmt.Println(dl[1],dl[2])
	}
	pReg := regexp.MustCompile(`<p>([\s\S]*)</p>`).FindAllStringSubmatch(str,-1)
	if len(pReg) > 0{
		p := pReg[0][1]
		p = strings.Replace(p," ","",-1)
		fmt.Println(p)
	}

}
