package jdSpider

import (
	"../sqlgo"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	eachGood = `<li.*class="gl-item">[\s\S]*?</div>[.\s]*</li>`
	goodPrice = `<em>￥</em><i>(.*)</i>`
	goodHref = `<div class="p-name p-name-type-3">[\s\S]*?href="(.*?)"[\s\S]*?</div>`
	goodName = `<div class="p-name p-name-type-3">[\s\S]*<em>(<.*>)?([\S\s]*)</em>[\n.\s\S]*</div>`
	goodPtable = `<div class="Ptable-item">[\s\S]*?</div>`
	goodPackage = `<div class="package-list">[\s\S]*?</div>`
	goodParameter = `<ul class="parameter2 p-parameter-list">[\s\S]*?</ul>`
)
const (
	getGidGhref = `select goods_id,goods_href from goodsTable`
)
/*
func main() {

	href := "https://list.jd.com/list.html?cat=1672,2577,3997"
	getGood(href)
	getGoodDetal

}


 */

func GetGood(href string,pid int,pwg *sync.WaitGroup,mu *sync.Mutex){
	defer pwg.Done()
	//getClassDetal(href,pid,mu)
	getClassDetal(href,pid,mu)
}

func GetGoodsIntr(){


	result := sqlgo.SelectSql(getGidGhref)
	fmt.Println(len(result))
	for k,v:= range result{
		time.Sleep(100*time.Millisecond)
		introJson,_ := json.Marshal(GetGoodDetal(v[1]))
		sqlgo.InsertGoodIntro(v[0],string(introJson))
		fmt.Println("id :",k)

	}


}


func getClassDetal(href string,pid int,mu *sync.Mutex) []string{
	//href := "https://list.jd.com/list.html?cat=9987,830,13661"
	hrefList := []string{}
	defer func() {
		if len(hrefList) > 0{
			fmt.Println("end: ",len(hrefList))
		}
	}()

	resp,err:=http.Get(href)
	if err != nil {
		ctopen := fmt.Sprintf("can not open page %s\n",href)
		mu.Lock()
		file,_ := os.OpenFile("./out.txt",os.O_WRONLY|os.O_CREATE|os.O_APPEND,0666)
		io.WriteString(file,ctopen)
		mu.Unlock()

		return hrefList
	}
	bBody,err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ctread:=fmt.Sprintf("can not ReadAll page %s\n",href)
		mu.Lock()
		file,_ := os.OpenFile("./out.txt",os.O_WRONLY|os.O_CREATE|os.O_APPEND,0666)
		io.WriteString(file,ctread)
		mu.Unlock()
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

			sqlgo.InsertGood(pid,itemName,itemPriceFloat,itemHref)
			hrefList = append(hrefList,itemHref)
		}
	}
	if len(goodList) == 0{
			ctmatch := fmt.Sprintf("can NOT match %s \n",href)
			mu.Lock()
			file,_ := os.OpenFile("./out.txt",os.O_WRONLY|os.O_CREATE|os.O_APPEND,0666)
			io.WriteString(file,ctmatch)
			mu.Unlock()
	}
	return hrefList
}
func GetGoodDetal(href string) map[string]interface{} {
	var information map[string]interface{}
	information = make(map[string]interface{})
	if href[:6] != "https:"{
		href =  "https:"+href
	}
	resp,err:=http.Get(href)
	if err != nil {
		fmt.Printf("can't open page %s\n",href)
		return information
	}
	bBody,err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Printf("can't Read page %s\n",href)
		return information
	}

	body := string(bBody)
	resp.Body.Close()

	parameterList := regexp.MustCompile(goodParameter).FindAllString(body,-1)
	if len(parameterList) > 0 {
		paserIntro(parameterList[0],information)
	}

	ptableList := regexp.MustCompile(goodPtable).FindAllString(body,-1)
	for _,item := range ptableList{
		paserPackage(item,information)
	}
	packageList := regexp.MustCompile(goodPackage).FindAllString(body,-1)
	if len(packageList) > 0 {
		paserPackage(packageList[0],information)
	}
	//fmt.Println(information)
	return information
}

func paserIntro(str string,info map[string]interface{}){
	//fmt.Println("title:商品介绍")
	liReg := regexp.MustCompile(`<li title=[\s\S]*?>(.*)</li>`).FindAllStringSubmatch(str,-1)
	for _,li := range(liReg){
		kv := strings.Split(li[1],"：")
		havA,_ := regexp.MatchString(`<a[\s\S]*?</a>`,kv[1])
		if havA == true{
			result := regexp.MustCompile(`<a[\s\S]*?>(.*)</a>`).FindAllStringSubmatch(kv[1],1)
			info[kv[0]] = result[0][1]
		}else{
			info[kv[0]] = kv[1]
		}

	}
}

func paserPackage(str string,info map[string]interface{}){
	//fmt.Println("=======paserItem==========")
	h3Reg := regexp.MustCompile(`<h3>(.*)</h3>`).FindAllStringSubmatch(str,-1)
	title := h3Reg[0][1]

	//fmt.Println("title:",title)
	dlReg := regexp.MustCompile(`<dt>(.*)</dt><dd>(.*)</dd>`).FindAllStringSubmatch(str,-1)
	if len(dlReg) > 0 {
		var dlmap map[string]string
		dlmap = make(map[string]string)
		for _,dl := range(dlReg){
			//fmt.Println(dl[1],dl[2])
			dlmap[dl[1]] = dl[2]
		}
		info[title]=dlmap
	}

	pReg := regexp.MustCompile(`<p>([\s\S]*)</p>`).FindAllStringSubmatch(str,-1)
	if len(pReg) > 0{
		p := pReg[0][1]
		p = strings.Replace(p," ","",-1)
		info["包装清单"]=p
	}
}

