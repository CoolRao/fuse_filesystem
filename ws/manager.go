package ws

import (
	"fuse_file_system/log"
)

var ClientManager = NewManger()

type Manger struct {
	Register   chan *Client
	UnRegister chan *Client
	Clients    map[string]*Client
	exit       chan string
}

func NewManger() *Manger {
	return &Manger{
		Register:   make(chan *Client, 100),
		UnRegister: make(chan *Client, 100),
		Clients:    make(map[string]*Client),
		exit:       make(chan string),
	}
}

func (m *Manger) Close() {
	for ip, client := range m.Clients {
		log.Logger.Debugf("client exit %s ", ip)
		client.close()
	}
}

func (m *Manger) Client(ip string) *Client {
	if client, ok := m.Clients[ip]; ok {
		return client
	}
	return nil
}

func (m *Manger) Run() {
	for {
		select {
		case exit := <-m.exit:
			log.Logger.Warnf("ws manager exit %s ", exit)
			return
		case client := <-m.Register:
			log.Logger.Infof("ws client registered  %s ", client.Ip)
			m.Clients[client.Ip] = client
		case client := <-m.UnRegister:
			log.Logger.Warnf("ws client unregistered %s ", client.Ip)
			if _, ok := m.Clients[client.Ip]; ok {
				delete(m.Clients, client.Ip)
			}
		}
	}
}
