package main

import (
	"fuse_file_system/log"
	"fuse_file_system/ws"
	lg "log"
)

func main() {
	log.InitLogger()
	go ws.ClientManger().Run()
	defer ws.ClientManger().Close()
	err := ws.NewWebSocketServer("0.0.0.0:8888")
	if err != nil {
		lg.Fatal(err.Error())
	}
}
