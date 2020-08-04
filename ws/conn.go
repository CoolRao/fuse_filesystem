package ws

import (
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
		exitR:            make(chan string),
		exitW:            make(chan string),
		exitH:            make(chan string),
		intervalBeatTime: 5 * 60,
	}
}

func (c *Client) read() {
	for {
		select {
		case exit := <-c.exitR:
			log.Logger.Warnf("c client exit read  %s", exit)
			return
		default:
			msgType, bytes, err := c.socket.ReadMessage()
			if err != nil {
				log.Logger.Errorf("c client read error %s ", err.Error())
				continue
			}
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
				log.Logger.Warnf("wc client recv close message %s ", string(bytes))
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
				log.Logger.Warnf("unknown message %d", msgType)
			}

		}
	}
}

func (c *Client) Read() chan RecvMsg {
	return c.recv
}

func (c *Client) Write() {

}

func (c *Client) isTimeOut() bool {
	if int64(time.Now().Sub(c.heartTime).Seconds()) > c.intervalBeatTime {
		return true
	}
	return false
}

func (c *Client) heartBeat() {
	for {
		select {
		case exit := <-c.exitH:
			log.Logger.Warnf("c client exit heart beat  %s ", exit)
			return
		default:
			time.Sleep(5 * time.Minute)
			if c.isTimeOut() {
				log.Logger.Warnf("heart time time out  ,start retry connect")
				continue
			} else {
				c.send <- SendMsg{Type: websocket.PingMessage, Body: []byte("hello")}
				log.Logger.Infof("send heart beat ")
			}
		}
	}
}

func (c *Client) write() {
	for {
		select {
		case exit := <-c.exitW:
			log.Logger.Warnf("c client exit write %s ", exit)
			return
		case msg := <-c.send:
			err := c.socket.WriteMessage(msg.Type, msg.Body)
			if err != nil {
				log.Logger.Errorf("c client write error %s ", err.Error())
				continue
			}
			log.Logger.Debugf("c client write : %d %s ", msg.Type, string(msg.Body))
		}
	}
}


func (c *Client) closeSocket() {
	if c.socket != nil {
		err := c.socket.Close()
		if err != nil {
			log.Logger.Errorf("c client socket close error %s ", err.Error())
		}
	}
}

func (c *Client) close() {
	c.closeSocket()
	c.exitR <- "exit read"
	c.exitW <- "exit write"
	c.exitH <- "exit heart beat"
	clientManager.UnRegister <- c
}
