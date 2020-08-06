package ws

import (
	"fuse_file_system/config"
	"fuse_file_system/log"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)



func NewWebSocketServer(addr string) error {
	gin.SetMode(gin.ReleaseMode)
	httpRouter := gin.New()
	httpRouter.GET(config.ConnPath, WebSocketConnHandler)
	httpRouter.Use(gin.Recovery())
	httpSever := &http.Server{
		Addr:           addr,
		Handler:        httpRouter,
		ReadTimeout:    15 * time.Second,
		WriteTimeout:   15 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	err := httpSever.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}


func WebSocketConnHandler(c *gin.Context) {
	log.Logger.Infof("client connect %s \n", c.ClientIP())
	conn, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		http.NotFound(c.Writer, c.Request)
		return
	}
	client := NewClient(c.ClientIP(), conn)
	ClientManager.Register <- client
	go client.Read()
	go client.Write()
	go client.HeartBeat()
}
