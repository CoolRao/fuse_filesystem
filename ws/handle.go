package ws

import (
	"fuse_file_system/log"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

func WebSocketConnHandler(c *gin.Context) {
	log.Logger.Infof("client connect %s \n", c.ClientIP())
	conn, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		http.NotFound(c.Writer, c.Request)
		return
	}
	client := NewClient(c.ClientIP(), conn)
	clientManager.Register <- client
	go client.Read()
	go client.Write()
	go client.heartBeat()
}
