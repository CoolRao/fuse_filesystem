package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"fuse_file_system/log"
	"fuse_file_system/model"
	"fuse_file_system/ws"
	"github.com/gorilla/websocket"
	"os"
)

type FileHandler struct {
	client  *ws.WsClient
	workDir string
	host    string
	exit    chan string
}

func NewFileHandler(host string, workDir string) *FileHandler {
	client := ws.NewWsClient(host)
	return &FileHandler{
		client:  client,
		host:    host,
		workDir: workDir,
		exit:    make(chan string),
	}
}

func (fh *FileHandler) Close() {
	fh.client.Close()
	fh.exit <- "exit"
}

func (fh *FileHandler) sync() {
	for {
		select {
		case <-fh.exit:
			log.Logger.Warnln("fileHandler exit")
			return
		case bytes := <-fh.client.Read():
			msg := ws.Message{}
			err := json.Unmarshal(bytes.Body, &msg)
			if err != nil {
				log.Logger.Errorf("fileHandler: read json unmarshal error %s ", err.Error())
				continue
			}
			switch msg.MsgType {
			case ws.FileSyncAttr:
				result, err := fh.SyncFileAttr(msg.DirType)
				if err != nil {
					log.Logger.Errorf("file state fail %s \n", err.Error())
				}
				fh.send(msg, result)
				break
			case ws.FileReadType:
				log.Logger.Warnln(string(bytes.Body))
				result, err := fh.FileRead(msg.FileName, msg.Size, msg.Off)
				if err != nil {
					log.Logger.Errorf("file read fail %s \n", err.Error())
				}
				fh.send(msg, result)
				break
			case ws.FileStateType:
				result, err := fh.FileState(msg.FileName)
				if err != nil {
					log.Logger.Errorf("file state fail %s \n", err.Error())
				}
				fh.send(msg, result)
				break
			default:
				log.Logger.Warnln("not support msgType")
			}
		}
	}
}

func (fh *FileHandler) send(msg ws.Message, result []byte) {
	msg.Body = result
	content, err := json.Marshal(msg)
	if err != nil {
		log.Logger.Errorf("json fail  %s ", err.Error())
	}
	fh.client.Write(ws.SendMsg{Type: websocket.TextMessage, Body: content})
}

func (fh *FileHandler) run() error {

	go fh.sync()

	defer fh.Close()

	err := fh.client.Run()
	if err != nil {
		return err
	}
	return nil
}

func (fh *FileHandler) FileRead(fileName string, size, off int64) ([]byte, error) {
	absolutePath, ok, err := GetAbsolutePath(fh.workDir, fileName)

	if err != nil {
		log.Logger.Errorln("fileRead: get file path  error ", err.Error())
		return nil, err
	}
	if !ok {
		log.Logger.Infoln("fileRead: file path no found in work storage ")
		return nil, fmt.Errorf("no find file %s  \n", fileName)
	}

	file, err := os.Open(absolutePath)
	if err != nil {
		log.Logger.Errorln("fileRead: open file error ", absolutePath, err.Error())
		return nil, err
	}

	defer func() {
		if file != nil {
			file.Close()
		}
	}()
	bufReader := bufio.NewReader(file)
	_, err = file.Seek(off, 0)
	if err != nil {
		log.Logger.Errorln("fileRead: seek file error  ", absolutePath, err.Error())
		return nil, err
	}

	bytes := make([]byte, size)
	_, err = bufReader.Read(bytes)
	if err != nil {
		log.Logger.Errorln("fileRead: read file []byte error  ", absolutePath, err.Error())
		return nil, err
	}
	log.Logger.Infoln("FileRead: ", fileName, size, off)
	return bytes, nil
}

func (fh *FileHandler) SyncFileAttr(dirType string) ([]byte, error) {
	dir, err := TraverseDir(fh.workDir, dirType)
	if err != nil {
		return nil, err
	}
	bytes, err := json.Marshal(dir)
	return bytes, err
}

func (fh *FileHandler) FileState(fileName string) ([]byte, error) {
	absolutePath, ok, err := GetAbsolutePath(fh.workDir, fileName)
	if err != nil {
		log.Logger.Errorln("fileState: get file path  error ", err.Error(), fileName)
		return nil, err
	}
	if !ok {
		log.Logger.Infoln("fileState: file path no found in work storage ")
		return nil, fmt.Errorf("no find file %s", fileName)
	}
	fileInfo, err := os.Stat(absolutePath)
	if err != nil {
		log.Logger.Errorln("fileState: get file state   error ", err.Error())
		return nil, err
	}
	log.Logger.Infoln("fileState:  ", fileName)
	bytes, err := json.Marshal(model.CovertFileState(fileInfo))
	return bytes, err
}
