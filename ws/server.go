package ws

import (
	"fuse_file_system/config"
	"fuse_file_system/log"
	"fuse_file_system/utils"
	"github.com/gorilla/websocket"
	"time"
)

type Client struct {
	socket           *websocket.Conn
	send             chan SendMsg
	recv             chan RecvMsg
	Ip               string
	Id               string
	exitR            chan string // 读退出
	exitW            chan string // 写退出
	exitH            chan string // 心跳退出
	intervalBeatTime int64       // 心跳间隔
	heartTime        time.Time   // 最近一次心跳时间
}

func NewClient(Ip string, conn *websocket.Conn) *Client {
	return &Client{
		Ip:               Ip,
		socket:           conn,
		Id:               utils.GetUUID(),
		send:             make(chan SendMsg, 100),
		recv:             make(chan RecvMsg, 100),
		exitR:            make(chan string,1),
		exitW:            make(chan string,1),
		exitH:            make(chan string,1),
		intervalBeatTime: 5 * 60,
	}
}

func (c *Client) Read() {
	for {
		select {
		case exit := <-c.exitR:
			log.Logger.Warnf("c client exit Read  %s", exit)
			return
		default:
			msgType, bytes, err := c.socket.ReadMessage()
			if err != nil {
				log.Logger.Errorf("c client Read error %s ", err.Error())
				c.close()
				return
			}
			//log.Logger.Debugf("c client receive  %d  %s ", msgType, string(bytes))
			switch msgType {
			case websocket.PingMessage:
				log.Logger.Infof("c client recv ping message %s ", string(bytes))
				c.send <- SendMsg{Type: websocket.PongMessage}
				break
			case websocket.PongMessage:
				log.Logger.Infof("c client recv pong message %s ", string(bytes))
				c.send <- SendMsg{Type: websocket.PingMessage}
				c.heartTime = time.Now()
				break
			case websocket.CloseMessage:
				log.Logger.Warnf("wc client recv Close message %s ", string(bytes))
				c.close()
				break
			case websocket.TextMessage:
				log.Logger.Debugf("c client recv text message %s ", string(bytes))
				c.recv <- RecvMsg{Type: websocket.TextMessage, Body: bytes}
				break
			case websocket.BinaryMessage:
				log.Logger.Debugf("c client recv binary message %d ", len(bytes))
				c.recv <- RecvMsg{Type: websocket.BinaryMessage, Body: bytes}
				break
			default:
				log.Logger.Warnf("unknown message %d  %s ", msgType, string(bytes))
			}

		}
	}
}

func (c *Client) Recv() chan RecvMsg {
	return c.recv
}

func (c *Client) Send(msg SendMsg) {
	c.send <- msg
}

func (c *Client) HeartBeat() {
	timer := time.NewTicker(config.ServerHeartTime * time.Second)
	for {
		select {
		case exit := <-c.exitH:
			log.Logger.Warnf("c client exit heart beat  %s ", exit)
			return
		case <-timer.C:
			err := c.socket.WriteMessage(websocket.PingMessage, []byte("hello"))
			if err != nil {
				log.Logger.Errorf("client ping message error %s ", err.Error())
				c.close()
				return
			}

		}
	}
}

func (c *Client) Write() {
	for {
		select {
		case exit := <-c.exitW:
			log.Logger.Warnf("c client exit Write %s ", exit)
			return
		case msg := <-c.send:
			err := c.socket.WriteMessage(msg.Type, msg.Body)
			if err != nil {
				log.Logger.Errorf("c client Write error %s ", err.Error())
				c.close()
				return
			}
			log.Logger.Debugf("c client Write : %d %s ", msg.Type, string(msg.Body))
		}
	}
}

func (c *Client) closeSocket() {
	if c.socket != nil {
		err := c.socket.Close()
		if err != nil {
			log.Logger.Errorf("c client socket Close error %s ", err.Error())
		}
	}
}

func (c *Client) close() {
	c.closeSocket()
	c.exitR <- "exit Read"
	c.exitW <- "exit Write"
	c.exitH <- "exit heart beat"
	clientManager.UnRegister <-c
}
