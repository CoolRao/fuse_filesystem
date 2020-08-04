package ws

import (
	"fuse_file_system/log"
	"github.com/gorilla/websocket"
	"net/url"
	"time"
)

type WsClient struct {
	socket           *websocket.Conn
	send             chan SendMsg
	recv             chan RecvMsg
	host             string
	exitR            chan string // 读退出
	exitW            chan string // 写退出
	exitH            chan string // 心跳退出
	intervalBeatTime int64       // 心跳间隔
	heartTime        time.Time   // 最近一次心跳时间
}

func NewWsClient(host string) *WsClient {
	return &WsClient{
		host:             host,
		send:             make(chan SendMsg, 100),
		recv:             make(chan RecvMsg, 100),
		exitR:            make(chan string),
		exitW:            make(chan string),
		exitH:            make(chan string),
		intervalBeatTime: 5 * 60,
	}
}

func (ws *WsClient) read() {
	for {
		select {
		case exit := <-ws.exitR:
			log.Logger.Warnf("ws client exit read  %s", exit)
			return
		default:
			msgType, bytes, err := ws.socket.ReadMessage()
			if err != nil {
				log.Logger.Errorf("ws client read error %s ", err.Error())
				ws.retryConn()
				continue
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
				log.Logger.Warnf("wc client recv close message %s ", string(bytes))
				ws.close()
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

func (ws *WsClient) isTimeOut() bool {
	if int64(time.Now().Sub(ws.heartTime).Seconds()) > ws.intervalBeatTime {
		return true
	}
	return false
}

func (ws *WsClient) heartBeat() {
	for {
		select {
		case exit := <-ws.exitH:
			log.Logger.Warnf("ws client exit heart beat  %s ", exit)
			return
		default:
			time.Sleep(5 * time.Minute)
			if ws.isTimeOut() {
				log.Logger.Warnf("heart time time out  ,start retry connect")
				ws.retryConn()
				continue
			} else {
				ws.send <- SendMsg{Type: websocket.PingMessage, Body: []byte("hello")}
				log.Logger.Infof("send heart beat ")
			}
		}
	}
}

func (ws *WsClient) write() {
	for {
		select {
		case exit := <-ws.exitW:
			log.Logger.Warnf("ws client exit write %s ", exit)
			return
		case msg := <-ws.send:
			err := ws.socket.WriteMessage(msg.Type, msg.Body)
			if err != nil {
				log.Logger.Errorf("ws client write error %s ", err.Error())
				ws.retryConn()
				continue
			}
			log.Logger.Debugf("ws client write : type: %d  body:  %s ", msg.Type, string(msg.Body))
		}
	}
}

func (ws *WsClient) retryConn() {
	time.Sleep(3 * time.Second)
	log.Logger.Warnf("ws client retry connect %s ", ws.host)
	err := ws.dial()
	if err != nil {
		log.Logger.Errorf("ws client connect %s error %s ", ws.host, err.Error())
	}
}

func (ws *WsClient) closeSocket() {
	if ws.socket != nil {
		err := ws.socket.Close()
		if err != nil {
			log.Logger.Errorf("ws client socket close error %s ", err.Error())
		}
	}
}

func (ws *WsClient) close() {
	log.Logger.Warnf("ws client close ")
	ws.closeSocket()
	ws.exitR <- "exit read"
	ws.exitW <- "exit write"
	ws.exitH <- "exit heart beat"
}

func (ws *WsClient) Run() error {
	err := ws.dial()
	if err != nil {
		return err
	}
	defer ws.close()
	go ws.read()
	go ws.write()
	go ws.heartBeat()
	select {}

}

func (ws *WsClient) run() {

}

func (ws *WsClient) dial() error {
	u := url.URL{Scheme: "ws", Host: ws.host, Path: ConnPath}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return err
	}
	ws.socket = conn
	// 连接成功之后发心跳
	ws.Write(SendMsg{Type: websocket.PingMessage, Body: []byte("hello")})
	return nil
}
