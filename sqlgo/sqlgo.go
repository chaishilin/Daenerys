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

func Hello()  string {
	return "hello"
}

/*
func main() {


	//sqlCmd := "insert into small(id,sname,offset) values (4,?,0.8)"
	//curdSql(conn,sqlCmd,"life 2")

	conn :=InitMySql()
	sqlCmd := `
	create table %s(
	id int,
	mid int auto_increment,
	name varchar(20),
	primary key(mid))
	engine=InnoDB default charset=utf8
	`
	sqlCmd = fmt.Sprintf(sqlCmd,name)
	CurdSql(conn,sqlCmd)

	conn.Close()
	//sqlCmd := "select * from small where id = ? and sid <= ?"
	//selectSql(conn,sqlCmd,3, 7)



}

 */

func InitMySql() *sql.DB {
	dbInfo := fmt.Sprintf("%s:%s@%s(%s:%d)/%s",userName,passWord,netWork,server,port,dataBase)
	fmt.Println("log as : ",dbInfo)
	conn, err := sql.Open("mysql", dbInfo)
	if nil != err {
		fmt.Println("connect db error: ", err)
	}
	return conn
}

func selectSql(conn *sql.DB,sqlCmd string,args... interface{}){
	////查询所有
	rows, err := conn.Query(sqlCmd,args...)
	if err != nil {
		fmt.Println("Query error ",err.Error())
	}
	//查看所有列名
	cols, err:=rows.Columns()
	if err != nil {
		fmt.Println("Columns error ",err.Error())
	}
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
}



func CurdSql(conn *sql.DB,sqlCmd string,args... interface{}){

	stmt, err := conn.Prepare(sqlCmd)
	if err != nil {
		fmt.Println("prepare error ",err.Error())
	}
	_, err = stmt.Exec(args...)
	if err != nil {
		fmt.Println("can not exec arg",err.Error())
	}
	/*
	//获取最后插入的id
	rid, _ := rs.LastInsertId()
	fmt.Println("id:", rid)
	//获取影响的行数
	affectNum, _ := rs.RowsAffected()
	fmt.Println("affectNum:", affectNum)
	 */
}
