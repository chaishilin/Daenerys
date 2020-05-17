package jdSpider

import (
	"../sqlgo"
	"database/sql"
	"fmt"
	"regexp"
)

const (
	itemReg  = `<!--  /widget/cat-item/cat-item.vm -->[\s\S]+?<!--/ /widget/cat-item/cat-item.vm -->`
	titleReg = `<span>(.+)</span>`
	dlReg    = `<dl class=\"clearfix\"[\s\S]*?</dl>`
	dtReg    = `<dt>[\s\S]+?</dt>`
	ddReg    = `<dd>[\s\S]+?</dd>`
	infoReg  = `<a href="([\/\.\?\-,=&_0-9a-zA-Z]+)"[^>]*>([^<]+)</a>`
)

//select a.class_id,b.pid,a.class_name from classTable as a join classRelate as b where a.class_id = b.class_id and b.pid = 1108;

var count = 0

// 能不能把count改成管道？


//var conn *sql.DB
/*
func main() {
	url := "https://www.jd.com/allSort.aspx"
	resp, _ := http.Get(url)
	b, _ := ioutil.ReadAll(resp.Body)
	s := fmt.Sprintf("%s", b)
	resp.Body.Close()
	conn = sqlgo.InitMySql()
	sqlgo.CreateTables(conn)
	doSpider(s)

}

 */
func DoSpider(htmlMsg string,conn *sql.DB) {
	items := regexp.MustCompile(itemReg).FindAllString(htmlMsg, -1)

	for _, each := range items {
		title := regexp.MustCompile(titleReg)
		titles := title.FindAllSubmatch([]byte(each), -1)

		count++
		fmt.Println("--", count, 0, string(titles[0][1]))
		sqlgo.InsertClass(conn,count,string(titles[0][1]),"")
		sqlgo.InsertRelate(conn,count,0)

		pid := count
		dls := regexp.MustCompile(dlReg).FindAllString(each, -1)

		for _, eachdl := range dls {
			func(eachdl string) {
				dtresults, pid := findMatch(eachdl, dtReg, pid)
				for _, v := range dtresults {
					count++
					fmt.Println("----", count, pid, v[1], v[0])
					sqlgo.InsertClass(conn,count,v[1],v[0])
					sqlgo.InsertRelate(conn,count,pid)
				}
				dlresults, pid := findMatch(eachdl, ddReg, count)
				for _, v := range dlresults {
					count++
					fmt.Println("--------", count, pid, v[1], v[0])
					sqlgo.InsertClass(conn,count,v[1],v[0])
					sqlgo.InsertRelate(conn,count,pid)
				}
			}(eachdl)
		}
	}
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


