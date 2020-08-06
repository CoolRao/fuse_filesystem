package main

import (
	"fuse_file_system/log"
	"fuse_file_system/ws"
	"github.com/gorilla/websocket"
	lg "log"
	"time"
)

func main() {
	log.InitLogger()
	go ws.NewWebSocketServer("0.0.0.0:8888")
	go ws.GetClientManager().Run()
	timer := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-timer.C:
			lg.Println("timer ....")
			client := ws.GetClientManager().Client("127.0.0.1")
			if client != nil {
				lg.Println("send message ")
				client.Send(ws.SendMsg{Type: websocket.TextMessage, Body: []byte("hello ......")})
			}
		}
	}
	select {}
}
