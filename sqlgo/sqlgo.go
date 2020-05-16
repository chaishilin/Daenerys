package sqlgo

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
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
	fmt.Println("log as : ",dbInfo)
	conn, err := sql.Open("mysql", dbInfo)
	if nil != err {
		fmt.Println("connect db error: ", err)
	}
	return conn
}

func CreateTables(conn *sql.DB){
	sqlCmd := fmt.Sprintf(`
								create table IF NOT EXISTS %s(
								class_id int,
								pid int,
								primary key(class_id))
								engine=InnoDB default charset=utf8
								`,"classRelate")

	CurdSql(conn, sqlCmd)


	sqlCmd = fmt.Sprintf(`
								create table IF NOT EXISTS %s(
								class_id int auto_increment,
								class_name varchar(20),
								class_href varchar(50),
								primary key(class_id))
								engine=InnoDB default charset=utf8
								`,"classTable")
	CurdSql(conn, sqlCmd)


	sqlCmd = fmt.Sprintf(`
								create table IF NOT EXISTS %s(
								goods_id int auto_increment,
								class_id int,
								goods_name varchar(20),
								goods_price float,
								goods_href varchar(50),
								primary key(goods_id))
								engine=InnoDB default charset=utf8
								`,"classTable")
	CurdSql(conn, sqlCmd)
}

func InsertRelate(conn *sql.DB,class_id int,pid int){

	sqlCmd := fmt.Sprintf(`insert into classRelate
									(class_id,pid)
									values
									(%d,%d)
									`, class_id,pid)
	CurdSql(conn, sqlCmd)
}
func InsertClass(conn *sql.DB,class_id int,class_name string,class_href string){
	sqlCmd := fmt.Sprintf(`insert into classTable
									(class_id,class_name,class_href)
									values
									(%d,'%s','%s')
									`, class_id,class_name,class_href)
	CurdSql(conn, sqlCmd)
}

func SelectSql(conn *sql.DB,sqlCmd string,args... interface{}) string{
	////查询所有
	rows, err := conn.Query(sqlCmd,args...)
	if err != nil {
		fmt.Println("Query error ",err.Error())
		return ""
	}
	//查看所有列名
	cols, err:=rows.Columns()
	if err != nil {
		fmt.Println("Columns error ",err.Error())
		return ""
	}
	msg := ""

	vals := make([]sql.RawBytes, len(cols))
	//vals转换为interface, 查看https://github.com/golang/go/wiki/InterfaceSlice
	scanArgs := make([]interface{}, len(vals))
	for i := range vals {
		scanArgs[i] = &vals[i]
	}
	for rows.Next() {
		rows.Scan(scanArgs...)
		for _, val := range vals{
			msg = fmt.Sprintf("%s,%s",msg,string(val))
		}
		//println()
	}
	return msg
}