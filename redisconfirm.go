package main

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
)

type redisState int32

const (
	getOk redisState = iota
	getErr
	setOk
	setExist
	setErr
)

func main() {

	conn := InitRedis()
	register, code, _ := Register(&conn, "cn", "pss")
	fmt.Println(register,stateParser(code))
	conn.Close()

}

func InitRedis() redis.Conn {
	conn, err := redis.Dial("tcp", ":6379")
	if err != nil {
		fmt.Println("redis connect error : ", err)
		return nil
	}
	return conn
}

func GetUsername(conn *redis.Conn, name string) (string, redisState, error) {
	resp, err := (*conn).Do("get", name)
	if err != nil {
		err = fmt.Errorf("redis set error : ", err)
		return "", getErr, err
	}
	resMsg := fmt.Sprintf("%v", resp)
	if resMsg == "<nil>" {
		return resMsg, getErr, nil
	}
	resMsg = fmt.Sprintf("%s", resp)
	return resMsg, getOk, nil
}

func Register(conn *redis.Conn, name string, password string) (string, redisState, error) {
	var resMsg string

	respGet, getState, _ := GetUsername(conn, name)

	if getState == getOk {
		return respGet, setExist, nil
	}

	respDo, err := (*conn).Do("set", name, password)
	if err != nil {
		err = fmt.Errorf("redis set error : %s ", err)
		return "", setErr, err
	}
	resMsg = fmt.Sprintf("%s", respDo)
	return resMsg, setOk, nil
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
