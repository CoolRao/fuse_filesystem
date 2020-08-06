package fusecore

import (
	"fuse_file_system/log"
	"fuse_file_system/ws"
	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
	"github.com/robfig/cron/v3"
)

var copyDir string

var ClientManager *ws.Manger


type FuseManage struct {
	server  *fuse.Server
	root    *Root
	dirType int64
	timer   *cron.Cron
	timerId cron.EntryID
	host    string
}

func NewFuseManage(mountPath, dest string, host string) (*FuseManage, error) {
	opts := &fs.Options{}
	opts.Debug = false
	opts.Options = []string{"nonempty"}
	copyDir = dest
	root := NewRoot(mountPath, dest)
	server, err := fs.Mount(mountPath, root, opts)
	if err != nil {
		return nil, err
	}
	return &FuseManage{
		server: server,
		root:   root,
		host:   host,
		timer:  cron.New(cron.WithSeconds()),
	}, nil
}

func (fm *FuseManage) Run() error {
	ClientManager = ws.NewManger()
	log.Logger.Warnln(ClientManager)
	go NewWebSocketServer(fm.host)
	go ClientManager.Run()
	fm.server.Wait()
	return nil
}

func (fm *FuseManage) Close() error {

	return nil
}
