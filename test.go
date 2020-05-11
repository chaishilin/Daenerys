package main

import (
	"./sqlgo"
	"fmt"
)

func main() {


	//sqlCmd := "insert into small(id,sname,offset) values (4,?,0.8)"
	//curdSql(conn,sqlCmd,"life 2")
	conn := sqlgo.InitMySql()

	sqlCmd := fmt.Sprintf(`
								create table IF NOT EXISTS %s(
								bid int auto_increment,
								name varchar(100),
								href varchar(100),
								primary key(bid))
								engine=InnoDB default charset=utf8
								`,"总目录")
	sqlgo.CurdSql(conn, sqlCmd)
	sqlCmd = fmt.Sprintf(`insert into 总目录
									(name)
									values
									('%s')
									`, `h`)
	fmt.Println(sqlCmd)
	sqlgo.CurdSql(conn, sqlCmd)

	/*
	sqlCmd := `
	create table %s(
	id int,
	mid int auto_increment,
	name varchar(20),
	primary key(mid))
	engine=InnoDB default charset=utf8
	`
	sqlCmd = fmt.Sprintf(sqlCmd,"chhh")

	 */
	//sqlgo.CurdSql(conn,sqlCmd)

	//conn.Close()
	//sqlCmd := "select * from small where id = ? and sid <= ?"
	//selectSql(conn,sqlCmd,3, 7)



}