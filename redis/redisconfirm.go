package redisconfirm

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
)

type redisState int32

const (
	GetOk redisState = iota
	GetErr
	SetOk
	SetExist
	SetErr
)
/*
func main() {

	conn := InitRedis()
	register, code, _ := Register(&conn, "cn", "pss")
	fmt.Println(register,stateParser(code))
	conn.Close()

}
*/

func InitRedis() redis.Conn {
	conn, err := redis.Dial("tcp", ":6379")
	if err != nil {
		fmt.Println("redis connect error : ", err)
		return nil
	}
	return conn
}





func getUsername(conn *redis.Conn, name string) (string, redisState, error) {
	resp, err := (*conn).Do("get", name)
	if err != nil {
		err = fmt.Errorf("redis set error : ", err)
		return "", GetErr, err
	}
	resMsg := fmt.Sprintf("%v", resp)
	if resMsg == "<nil>" {
		return resMsg, GetErr, nil
	}
	resMsg = fmt.Sprintf("%s", resp)
	return resMsg, GetOk, nil
}
func LogCheck (conn *redis.Conn,name string, password string) bool {

	respGet, getState, _ :=  getUsername(conn,name)
	if getState != GetOk {
		return false
	}else{
		if respGet == password{
			return true
		}else{
			return false
		}
	}
}
func Register(conn *redis.Conn, name string, password string) (string, redisState, error) {
	var resMsg string

	respGet, getState, _ := getUsername(conn, name)

	if getState == GetOk {
		return respGet, SetExist, nil
	}

	respDo, err := (*conn).Do("set", name, password)
	if err != nil {
		err = fmt.Errorf("redis set error : %s ", err)
		return "", SetErr, err
	}
	resMsg = fmt.Sprintf("%s", respDo)
	return resMsg, SetOk, nil
}

func stateParser(code redisState) string {
	codeInfo := ""
	switch code {
	case 0:
		codeInfo = "getOk"
	case 1:
		codeInfo = "getErr"
	case 2:
		codeInfo = "setOk"
	case 3:
		codeInfo = "setExist"
	case 4:
		codeInfo = "setErr"
	}
	return codeInfo
}
