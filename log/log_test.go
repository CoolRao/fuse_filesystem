package log

import (
	"fmt"
	"testing"
)

func TestLog(t *testing.T){
	err := InitLogger("./log", "log.log")
	if err!=nil{
		fmt.Println(err.Error())
		return
	}
	Logger.Infoln("hello world")

}
