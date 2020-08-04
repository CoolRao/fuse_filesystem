package main

import (
	"fmt"
	"fuse_file_system/log"
)

func main() {
	err := log.InitLogger("./log", "slave.log")
	if err!=nil{
		fmt.Println(err.Error())
	}
	log.Logger.Infof(" %s","xjrw ")
	log.Logger.Infof(" %s \n","xjrw ")
}
