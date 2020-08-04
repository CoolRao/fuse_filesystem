package ws

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

const ConnPath = "/v1/conn"

func NewWebSocketServer(addr string) error {
	gin.SetMode(gin.ReleaseMode)
	httpRouter := gin.New()
	httpRouter.GET(ConnPath, WebSocketConnHandler)
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
