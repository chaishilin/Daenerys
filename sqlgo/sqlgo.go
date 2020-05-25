package sqlgo

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
)

const(
	userName = "root"
	passWord = "csl"
	netWork = "tcp"
	server = "localhost"
	port = 3306
	dataBase = "jddb"
)

func InitMySql() *sql.DB {
	dbInfo := fmt.Sprintf("%s:%s@%s(%s:%d)/%s",userName,passWord,netWork,server,port,dataBase)
	//fmt.Println("log as : ",dbInfo)
	conn, err := sql.Open("mysql", dbInfo)
	conn.SetMaxOpenConns(10000)
	conn.SetMaxIdleConns(5000)
	if nil != err {
		fmt.Println("connect db error: ", err)
	}
	return conn
}

func DelAll(){
	CurdSql("drop table classRelate")
	CurdSql("drop table goodsTable")
	CurdSql("drop table classTable")
	CurdSql("drop table userInfo")
	CurdSql("drop table goodsIntro")
}

func CreateTables(){
	sqlCmd := fmt.Sprintf(`
								create table IF NOT EXISTS %s(
								class_id int,
								pid int,
								primary key(class_id))
								engine=InnoDB default charset=utf8
								`,"classRelate")

	CurdSql(sqlCmd)


	sqlCmd = fmt.Sprintf(`
								create table IF NOT EXISTS %s(
								class_id int auto_increment,
								class_name varchar(20),
								class_href text,
								primary key(class_id))
								engine=InnoDB default charset=utf8
								`,"classTable")
	CurdSql(sqlCmd)


	sqlCmd = fmt.Sprintf(`
								create table IF NOT EXISTS %s(
								goods_id int auto_increment,
								class_id int,
								goods_name varchar(100),
								goods_price float,
								goods_href text,
								primary key(goods_id))
								engine=InnoDB default charset=utf8
								`,"goodsTable")


	CurdSql(sqlCmd)

	sqlCmd = fmt.Sprintf(`
								create table IF NOT EXISTS %s(
								user_id int auto_increment,
								user_name varchar(20),
								user_email varchar(50),
								primary key(user_id))
								engine=InnoDB default charset=utf8
								`,"userInfo")
	CurdSql(sqlCmd)

	sqlCmd = fmt.Sprintf(`
								create table IF NOT EXISTS %s(
								goods_id int,
								intro json,
								primary key(goods_id))
								engine=InnoDB default charset=utf8
								`,"goodsIntro")



	CurdSql(sqlCmd)
}

func InsertGood(class_id int,goods_name string,goods_price float64,goods_href string){

	sqlCmd := `insert into goodsTable
				(class_id,goods_name,goods_price,goods_href)
				values
				(?,?,?,?)`

	CurdSql(sqlCmd,class_id,goods_name,goods_price,goods_href)
}

func InsertRelate(class_id int,pid int){

	sqlCmd := fmt.Sprintf(`insert into classRelate
									(class_id,pid)
									values
									(%d,%d)
									`, class_id,pid)
	CurdSql(sqlCmd)
}
func InsertClass(class_id int,class_name string,class_href string){
	sqlCmd := fmt.Sprintf(`insert into classTable
									(class_id,class_name,class_href)
									values(%d,'%s','%s')`,
									class_id,class_name,class_href)
	CurdSql(sqlCmd)
}
func InsertUser(username string,email string){
	sqlCmd := fmt.Sprintf(`insert into userInfo
									(user_name,user_email)
									values
									('%s','%s')
									`,username,email)
	CurdSql(sqlCmd)
}

func InsertGoodIntro(gid string,intro string){
	gidInt,_ := strconv.ParseInt(gid,10,64)
	sqlCmd := fmt.Sprintf(`insert into goodsIntro
									(goods_id,intro)
									values
									(?,?)
									`)
	CurdSql(sqlCmd,gidInt,intro)
}

func SelectSql(sqlCmd string,args... interface{}) [][]string{
	////查询所有
	//SelectSql(conn,"select * from classTable")
	dbInfo := fmt.Sprintf("%s:%s@%s(%s:%d)/%s",userName,passWord,netWork,server,port,dataBase)
	conn, err := sql.Open("mysql", dbInfo)
	defer conn.Close()
	rows, err := conn.Query(sqlCmd,args...)
	if err != nil {
		fmt.Println("Query error ",err.Error())
		return nil
	}
	//查看所有列名
	cols, err:=rows.Columns()
	if err != nil {
		fmt.Println("Columns error ",err.Error())
		return nil
	}
	result := [][]string{}

	vals := make([]sql.RawBytes, len(cols))
	//vals转换为interface, 查看https://github.com/golang/go/wiki/InterfaceSlice
	scanArgs := make([]interface{}, len(vals))
	for i := range vals {
		scanArgs[i] = &vals[i]
	}
	for rows.Next() {
		msg := []string{}
		rows.Scan(scanArgs...)
		for _, val := range vals{
			msg =append(msg,string(val))

		}
		result = append(result,msg)
	}
	return result
}


func CurdSql(sqlCmd string,args... interface{}){
	//fmt.Println(conn)
	dbInfo := fmt.Sprintf("%s:%s@%s(%s:%d)/%s",userName,passWord,netWork,server,port,dataBase)
	conn, err := sql.Open("mysql", dbInfo)
	stmt, err := conn.Prepare(sqlCmd)
	if err != nil {
		fmt.Println("prepare error ",err.Error(),"sql cmd:",sqlCmd)
		return
	}
	_, err = stmt.Exec(args...)
	if err != nil {
		fmt.Println("can not exec arg",err.Error(),"sql cmd:",sqlCmd)
	}
	defer conn.Close()

	/*
		//获取最后插入的id
		rid, _ := rs.LastInsertId()
		fmt.Println("id:", rid)
		//获取影响的行数
		affectNum, _ := rs.RowsAffected()
		fmt.Println("affectNum:", affectNum)
	*/
}