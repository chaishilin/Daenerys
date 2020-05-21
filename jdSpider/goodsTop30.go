package jdSpider

import (
	"../sqlgo"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
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

/*
func main() {

	href := "https://list.jd.com/list.html?cat=1672,2577,3997"
	getGood(href)

}
*/

func GetGood(href string,pid int,pwg *sync.WaitGroup){
	defer pwg.Done()
	getClassDetal(href,pid)
	/*
	itemHref := getClassDetal(href,pid)
	count := 0

	for _,item := range itemHref{
		getGoodDetal(item)
		count ++
	}
	 */
}

func getClassDetal(href string,pid int) []string{
	//href := "https://list.jd.com/list.html?cat=9987,830,13661"
	hrefList := []string{}
	defer func() {
		if len(hrefList) > 0{
			fmt.Println("end: ",len(hrefList))
		}
	}()

	resp,err:=http.Get(href)
	if err != nil {
		fmt.Printf("can not open page %s\n",href)
		return hrefList
	}
	bBody,err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("can not ReadAll page %s\n",href)
		return hrefList
	}
	body := string(bBody)
	resp.Body.Close()
	//获得分类下的商品页面，得到商品链接、名称、价格
	//还需要进入具体商品，得到商品简介、规格

	goodList := regexp.MustCompile(eachGood).FindAllString(body,-1)
	for _,good:=range goodList{
		//fmt.Println("----------------------")

		name:=regexp.MustCompile(goodName)
		nameResult := name.FindAllSubmatch([]byte(good),-1)
		price:=regexp.MustCompile(goodPrice)
		priceResult := price.FindAllSubmatch([]byte(good),-1)
		href:=regexp.MustCompile(goodHref)
		hrefResult := href.FindAllSubmatch([]byte(good),-1)
		if len(nameResult) > 0 && len(priceResult) > 0 && len(hrefResult) > 0 {
			//如果布局符合一般模式
			itemName := string(nameResult[0][2])
			itemName = strings.Replace(itemName,"\n","",-1)

			itemPrice :=string(priceResult[0][1])
			itemPriceFloat,_ := strconv.ParseFloat(itemPrice,10)

			itemHref := string(hrefResult[0][1])
			//fmt.Println("name: ",itemName,"price: ",itemPriceFloat,"href: ", itemHref[:10])

			conn := sqlgo.InitMySql()

			sqlgo.InsertGood(conn,pid,itemName,itemPriceFloat,itemHref)
			hrefList = append(hrefList,itemHref)
		}
	}
	if len(goodList) == 0{
			fmt.Println("can NOT match",href)
	}
	return hrefList
}
func getGoodDetal(href string){
	//href := <- hrefChan

	if href[:6] != "https:"{
		href =  "https:"+href
	}
	resp,err:=http.Get(href)
	if err != nil {
		fmt.Printf("can't open page %s\n",href)
		return
	}
	bBody,err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Printf("can't Read page %s\n",href)
		return
	}

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
