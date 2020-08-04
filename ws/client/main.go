package main

import (
	"fuse_file_system/log"
	"fuse_file_system/ws"
	lg "log"
)


func main() {
	log.InitLogger()
	wsClient := ws.NewWsClient("0.0.0.0:8888")
	err := wsClient.Run()
	if err!=nil{
		lg.Fatal(err.Error())
	}
}
