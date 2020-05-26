package main

import (
	"./redis"
	"./sqlgo"
	"encoding/json"
	"fmt"
	"github.com/dchest/captcha"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
	"sort"
	"strconv"
	"unicode/utf8"

)

type goodsInfo struct {
	Gname  string
	Gid    int
	Ghref  string
	Gprice float64
	Gjson string
}
type classInfo struct {
	Name      string
	Id        int
	Href      string
	GoodsList []goodsInfo
}
type classInfoList []classInfo

func (cList classInfoList) Len() int {return len(cList)}
func (cList classInfoList) Swap(i,j int) {cList[i],cList[j] = cList[j],cList[i]}
func (cList classInfoList)  Less(i,j int) bool {return len(cList[i].GoodsList) < len(cList[j].GoodsList)}

const (
	queryClassById = `select a.class_id,a.class_name,a.class_href 
					from classTable as a join classRelate as b 
					where a.class_id = b.class_id and 
					(b.pid = ? or b.class_id = ?);`
	queryClassByName = `select a.class_id,a.class_name,a.class_href 
						from classTable as a join classRelate as b 
						where a.class_id = b.class_id and 
						a.class_name regexp '%s';`
	SelectGoodsbyClassId = `select a.goods_id,a.goods_name,a.goods_href,a.goods_price 
							from goodsTable as a 
							where a.class_id = ?`
	queryGoodsById = `select a.goods_id,a.goods_name,a.goods_href,a.goods_price,b.intro
							from goodsTable as a left join goodsIntro as b
							on a.goods_id = b.goods_id 
							where a.goods_id = ?`
	queryGoodsByName = `select a.goods_id,a.goods_name,a.goods_href,a.goods_price,b.intro
							from goodsTable as a left join goodsIntro as b
							on a.goods_id = b.goods_id 
							where a.goods_name regexp '%s'`
)

var passwd string

func main() {
	passwd = os.Args[1:][0]
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", classHandler)
	mux.HandleFunc("/", logHandler)
	mux.HandleFunc("/regist", registHandler)
	mux.HandleFunc("/process", captchaVerify)
	mux.HandleFunc("/captcha/newId", newCapId)
	mux.HandleFunc("/template", makeTemplate)
	mux.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("./root/static"))))

	mux.Handle("/captcha/", captcha.Server(captcha.StdWidth, captcha.StdHeight))
	server := &http.Server{
		Addr:    ":18080",
		Handler: mux,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func newCapId(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, captcha.New())
}

func captchaVerify(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html;charset=utf-8")
	if !captcha.VerifyString(r.FormValue("captchaId"), r.FormValue("captchaSolution")) {
		//每个CaptchaId只能校验一次
		fmt.Fprint(w, "false")
	} else {
		fmt.Fprint(w, "true")
	}
}

func jsonPaser(str string) []string {

	var result []string
	var resultStr []string
	var resultMap []string
	var jsonMap map[string]interface{}
	jsonMap = make(map[string]interface{})
	json.Unmarshal([]byte(str),&jsonMap)

	for k, v := range jsonMap {
		switch fmt.Sprint(reflect.TypeOf(v)) {
		case "string":
			resultStr = append(resultStr,k+":	"+v.(string))
		case "map[string]interface {}":
			resultMap = append(resultMap,k)
			for key, value := range v.(map[string]interface{}) {
				resultMap = append(resultMap,key+":	"+value.(string))
			}
		}
	}
	result = append(resultStr,resultMap...)
	return result
}

func isTitle(str string) bool {
	for len(str) > 0 {
		r, size := utf8.DecodeRuneInString(str)
		if r == ':'{
			return false
		}
		str = str[size:]
	}
	return true
}

func makeTemplate(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	qType := r.Form.Get("type")
	query := r.Form.Get("input")
	if qType == "class"{
		_, err := strconv.Atoi(query)
		var result [][]string
		if err != nil {
			if len(query) > 0 {
				result = sqlgo.SelectSql(fmt.Sprintf(queryClassByName, fmt.Sprintf(".*%s.*", query))) //"'."+query+".'"
			}
		} else {
			result = sqlgo.SelectSql(queryClassById, query, query)
		}
		classList := makeClassInfo(result)
		t, err := template.ParseFiles("./root/classTemplate.html")

		if err != nil {
			w.WriteHeader(500)
			fmt.Fprint(w, err)
			return
		}
		t.Execute(w, classList)
	}else if qType == "goods"{
		_, err := strconv.Atoi(query)
		var result [][]string
		if err != nil {
			if len(query) > 0 {
				result = sqlgo.SelectSql(fmt.Sprintf(queryGoodsByName, fmt.Sprintf(".*%s.*", query))) //"'."+query+".'"
			}
		} else {
			result = sqlgo.SelectSql(queryGoodsById, query)
		}
		goodsList := makeGoodsInfo(result)
		//t, err := template.ParseFiles("./root/goodsTemplate.html")

		funcMap := template.FuncMap{
			"jPaser":jsonPaser,
			"isTitle":isTitle,
		}

		t := template.Must(template.New("goodsTemplate.html").Funcs(funcMap).ParseFiles("./root/goodsTemplate.html"))

		t.Execute(w, goodsList)
	}

}

func makeGoodsInfo(result [][]string) []goodsInfo{
	goodsList := []goodsInfo{}
	for _, k := range result {
		goodInfo := goodsInfo{}
		gId, _ := strconv.Atoi(k[0])
		goodInfo.Gid = gId
		goodInfo.Gname = k[1]
		goodInfo.Ghref = k[2]
		gPrice, _ := strconv.ParseFloat(k[3], 64)
		goodInfo.Gprice = gPrice
		goodInfo.Gjson = k[4]
		goodsList = append(goodsList, goodInfo)
	}
	return goodsList
}

func makeClassInfo(result [][]string) []classInfo {
	 var classList classInfoList

	for _, v := range result {
		classItem := classInfo{}
		cId, _ := strconv.Atoi(v[0])
		classItem.Id = cId
		classItem.Name = v[1]
		classItem.Href = v[2]

		goodItems := sqlgo.SelectSql(SelectGoodsbyClassId, cId)
		count := 0
		for _, k := range goodItems {
			goodInfo := goodsInfo{}
			gId, _ := strconv.Atoi(k[0])
			goodInfo.Gid = gId
			goodInfo.Gname = k[1]
			goodInfo.Ghref = k[2]
			gPrice, _ := strconv.ParseFloat(k[3], 64)
			goodInfo.Gprice = gPrice
			classItem.GoodsList = append(classItem.GoodsList, goodInfo)
			count ++
			if count >= 10{
				break
			}
		}

		classList = append(classList, classItem)
	}
	sort.Sort(sort.Reverse(classList))
	return classList
}



func classHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		file, _ := ioutil.ReadFile("./root/test2.html")
		fmt.Fprintf(w, string(file))
	} else if r.Method == "POST" {
		r.ParseForm()
		fmt.Fprint(w, 3)

	}

}

func logHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		file, _ := ioutil.ReadFile("./root/test.html")
		fmt.Fprintf(w, string(file))
	} else if r.Method == "POST" {
		r.ParseForm()
		conn := redisconfirm.InitRedis(passwd)
		logCheck := redisconfirm.LogCheck(&conn, r.Form.Get("name"), r.Form.Get("pwd"))
		fmt.Fprint(w, logCheck)
	}
}


func registHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		conn := redisconfirm.InitRedis(passwd)
		r.ParseForm()
		username := r.Form.Get("username")
		passwd := r.Form.Get("passwd")
		email_addr := r.Form.Get("email")
		_, state, _ := redisconfirm.Register(&conn, username, passwd)
		if state == redisconfirm.SetOk {
			fmt.Fprintf(w, "ok")
			//使用mysql存入用户邮箱、用户名，密码
			sqlgo.InsertUser(username, email_addr)
		} else {
			fmt.Fprintf(w, "exist")
		}
	}
}


