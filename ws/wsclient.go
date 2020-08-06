package ws

import (
	"fuse_file_system/config"
	"fuse_file_system/log"
	"github.com/gorilla/websocket"
	"net/url"
	"time"
)

type WsClient struct {
	socket     *websocket.Conn
	send       chan SendMsg
	recv       chan RecvMsg
	host       string
	exitR      chan string // 读退出
	exitW      chan string // 写退出
	exitH      chan string // 心跳退出
	heartTime  time.Time   // 最近一次心跳时间
	isRetrying bool        // 是否尝试重连中
}

func NewWsClient(host string) *WsClient {
	return &WsClient{
		host:  host,
		send:  make(chan SendMsg, 100),
		recv:  make(chan RecvMsg, 100),
		exitR: make(chan string,1),
		exitW: make(chan string,1),
		exitH: make(chan string,1),
	}
}

func (ws *WsClient) read() {
	for {
		select {
		case exit := <-ws.exitR:
			log.Logger.Warnf("ws client exit Read  %s", exit)
			return
		default:
			msgType, bytes, err := ws.socket.ReadMessage()
			if err != nil {
				log.Logger.Errorf("ws client Read error %s ", err.Error())
				ws.tryToConnectServer()
			}
			switch msgType {
			case websocket.PingMessage:
				log.Logger.Infof("ws client recv ping message %s ", string(bytes))
				ws.send <- SendMsg{Type: websocket.PongMessage, Body: []byte("pong")}
				break
			case websocket.PongMessage:
				log.Logger.Infof("ws client recv pong message %s ", string(bytes))
				ws.send <- SendMsg{Type: websocket.PingMessage, Body: []byte("ping")}
				ws.heartTime = time.Now()
				break
			case websocket.CloseMessage:
				log.Logger.Warnf("wc client recv Close message %s ", string(bytes))
				ws.Close()
				break
			case websocket.TextMessage:
				log.Logger.Debugf("ws client recv text message %s ", string(bytes))
				ws.recv <- RecvMsg{Type: websocket.TextMessage, Body: bytes}
				break
			case websocket.BinaryMessage:
				log.Logger.Debugf("ws client recv binary message %d ", len(bytes))
				ws.recv <- RecvMsg{Type: websocket.BinaryMessage, Body: bytes}
				break
			default:
				log.Logger.Warnf("unknown message %d", msgType)
			}

		}
	}
}

func (ws *WsClient) Read() chan RecvMsg {
	return ws.recv
}

func (ws *WsClient) Write(msg SendMsg) {
	ws.send <- msg
}

func (ws *WsClient) heartBeat() {
	ticker := time.NewTicker(config.ClientHeartTime * time.Second)
	for {
		select {
		case exit := <-ws.exitH:
			log.Logger.Warnf("ws client exit heart beat  %s ", exit)
			return
		case <-ticker.C:
			if ws.isRetrying {
				return
			}
			err := ws.socket.WriteMessage(websocket.PingMessage, []byte("hello"))
			if err != nil {
				log.Logger.Errorf("client ping message error %s ", err.Error())
				ws.tryToConnectServer()
			}
		}
	}
}

func (ws *WsClient) write() {
	for {
		select {
		case exit := <-ws.exitW:
			log.Logger.Warnf("ws client exit Write %s ", exit)
			return
		case msg := <-ws.send:
			err := ws.socket.WriteMessage(msg.Type, msg.Body)
			if err != nil {
				log.Logger.Errorf("ws client Write error %s ", err.Error())
				ws.tryToConnectServer()
			}
			log.Logger.Debugf("ws client Write : type: %d  body:  %s ", msg.Type, string(msg.Body))
		}
	}
}

func (ws *WsClient) tryToConnectServer() {
	if ws.isRetrying {
		return
	} else {
		ws.retryConn()
	}
}

//  todo 协程中调用
func (ws *WsClient) retryConn() {
	ws.isRetrying = true
	defer func() {
		ws.isRetrying = false
	}()
	log.Logger.Warnf("ws client retry connect %s ", ws.host)
	for {
		err := ws.dial()
		if err != nil {
			log.Logger.Errorf("ws client connect %s error %s ", ws.host, err.Error())
			time.Sleep(config.ClientRetryConnectTime * time.Second)
		} else {
			log.Logger.Infof("retry connect  server success")
			return
		}
	}

}

func (ws *WsClient) closeSocket() {
	if ws.socket != nil {
		err := ws.socket.Close()
		if err != nil {
			log.Logger.Errorf("ws client socket Close error %s ", err.Error())
		}
	}
}

func (ws *WsClient) Close() {
	log.Logger.Warnf("ws client Close ")
	ws.closeSocket()
	ws.exitR <- "exit Read"
	ws.exitW <- "exit Write"
	ws.exitH <- "exit heart beat"
}

func (ws *WsClient) Run() error {
	err := ws.dial()
	if err != nil {
		return err
	}
	defer ws.Close()
	ws.startWork()
	select {}

}

func (ws *WsClient) startWork() {
	go ws.read()
	go ws.write()
	go ws.heartBeat()
}

func (ws *WsClient) dial() error {
	u := url.URL{Scheme: "ws", Host: ws.host, Path: config.ConnPath}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return err
	}
	ws.socket = conn
	// 连接成功之后发心跳
	err = ws.socket.WriteMessage(websocket.PingMessage, []byte("hello"))
	if err != nil {
		return err
	}
	return nil
}
