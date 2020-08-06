package fusecore

import (
	"context"
	"encoding/json"
	"fmt"
	"fuse_file_system/config"
	"fuse_file_system/model"
	"fuse_file_system/utils"
	"fuse_file_system/ws"
	"github.com/gorilla/websocket"
	"sync"
	"time"
)

func getIpByFileName(fileName string) string {
	return "127.0.0.1"
}

func FileState(fileName string) (*model.FileStat, error) {
	ip := getIpByFileName(fileName)
	client := ClientManager.Client(ip)
	if client == nil {
		return nil, fmt.Errorf("client is nil %s  %s", ip, fileName)
	}
	msg := ws.Message{MsgType: config.FileStateType, FileName: fileName, MsgId: utils.GetUUID()}
	bytes, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	client.Send(ws.SendMsg{Type: websocket.TextMessage, Body: bytes})
	ctx, cancle := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancle()
	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("file state time out %s ", fileName)
		case result := <-client.Recv():
			message := ws.Message{}
			err := json.Unmarshal(result.Body, &message)
			if err != nil {
				return nil, err
			}
			if msg.MsgId == message.MsgId {
				state := &model.FileStat{}
				err := json.Unmarshal(message.Body, state)
				if err != nil {
					return nil, err
				}
				return state, nil
			}
		}
	}
}

func FileRead(fileName string, size, off int64) ([]byte, error) {
	ip := getIpByFileName(fileName)
	client := ClientManager.Client(ip)
	if client == nil {
		return nil, fmt.Errorf("client is nil %s  %s", ip, fileName)
	}
	sMsg := ws.Message{MsgType: config.FileReadType, FileName: fileName, MsgId: utils.GetUUID(), Size: size, Off: off}
	bytes, err := json.Marshal(sMsg)
	if err != nil {
		return nil, err
	}
	client.Send(ws.SendMsg{Type: websocket.TextMessage, Body: bytes})
	ctx, cancle := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancle()
	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("file read time out %s ", fileName)
		case result := <-client.Recv():
			rMsg := ws.Message{}
			err := json.Unmarshal(result.Body, &rMsg)
			if err != nil {
				return nil, err
			}
			if sMsg.MsgId == rMsg.MsgId {
				return rMsg.Body, nil
			}
		}
	}
}

func FileDirAttr(dirType string) ([]map[string]*model.FileStat, error) {
	sMsg := ws.Message{MsgType: config.FileSyncAttr, MsgId: utils.GetUUID(),DirType:dirType}
	bytes, err := json.Marshal(sMsg)
	if err != nil {
		return nil, err
	}
	var res []map[string]*model.FileStat
	wg := sync.WaitGroup{}
	for _, client := range ClientManager.Clients {
		wg.Add(1)
		go func() {
			client.Send(ws.SendMsg{Type: websocket.TextMessage, Body: bytes})
			ctx, _ := context.WithTimeout(context.TODO(), 10*time.Second)
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case result := <-client.Recv():
					message := ws.Message{}
					err := json.Unmarshal(result.Body, &message)
					if err != nil {
						return
					}
					if sMsg.MsgId == message.MsgId {
						attMap := make(map[string]*model.FileStat)
						err = json.Unmarshal(result.Body, &attMap)
						res = append(res, attMap)
						return
					}

				}
			}

		}()

	}
	wg.Wait()
	return res, nil
}




