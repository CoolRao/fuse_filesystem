package main

import (
	"fmt"
	"fuse_file_system/log"
	"fuse_file_system/ws"
	"github.com/gorilla/websocket"
	"time"
)

func main() {
	log.InitLogger()

	client := ws.NewWsClient("0.0.0.0:8888")
	go func() {
		timer := time.NewTicker(11 * time.Second)
		for {
			select {
			case <-timer.C:
				client.Write(ws.SendMsg{Type: websocket.TextMessage, Body: []byte("i`m coming")})
			}
		}
	}()
	err := client.Run()
	if err != nil {
		fmt.Println(err.Error())
	}
}
