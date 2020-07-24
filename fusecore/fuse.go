package fusecore

import (
	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
	"github.com/robfig/cron/v3"
)



type FuseManage struct {
	server  *fuse.Server
	root    *Root
	dirType int64
	timer   *cron.Cron
	timerId cron.EntryID
}

func NewFuseManage(mountPath string,) (*FuseManage, error) {
	opts := &fs.Options{}
	opts.Debug = false
	opts.Options=[]string{"nonempty"}
	root := NewRoot()
	server, err := fs.Mount(mountPath, root, opts)
	if err != nil {
		return nil, err
	}
	return &FuseManage{
		server:  server,
		root:    root,
		timer:   cron.New(cron.WithSeconds()),
	}, nil
}

func (fm *FuseManage) Run() error {
	fm.server.Wait()
	return nil
}

func (fm *FuseManage) Close() error {

	return nil
}

