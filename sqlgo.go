package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	//获取连接
	conn, err := sql.Open("mysql", "root:csl@tcp(localhost:3306)/tinydb")
	if nil != err {
		fmt.Println("connect db error: ", err)
	}

	//指定条件，指定参数查询
	/*
	rows, err := conn.Query("select * from user", 3)
	if nil != err{
		fmt.Println("query db error: ", err.Error())
		return
	}
	fmt.Println(rows)
	for rows.Next() {
		var name string
		//将值存入变量name中
		err= rows.Scan(&name)
		if err != nil{
			panic(err.Error())
		}
		fmt.Println(name)
	}

	 */

	////查询所有
	rows, _ := conn.Query("select * from user ")
	//查看所有列名
	cols, _:=rows.Columns()
	for _, col := range cols{
		println("col:", col)
	}
	vals := make([]sql.RawBytes, len(cols))
	//vals转换为interface, 查看https://github.com/golang/go/wiki/InterfaceSlice
	scanArgs := make([]interface{}, len(vals))
	for i := range vals {
		scanArgs[i] = &vals[i]
	}
	for rows.Next() {
		rows.Scan(scanArgs...)
		for _, val := range vals{
			print(string(val)," ")
		}
		println()
	}

	/*
	//保存数据(修改和删除操作)
	stmt, _ := conn.Prepare("insert into stu(id, no, name) values(?, ?, ?)")
	rs, _ := stmt.Exec(1, "003", "banana")
	//获取最后插入的id
	id, _ := rs.LastInsertId()
	fmt.Println("id:", id)
	//获取影响的行数
	affectNum, _ := rs.RowsAffected()
	fmt.Println("affectNum:", affectNum)


	 */
	conn.Close()
}