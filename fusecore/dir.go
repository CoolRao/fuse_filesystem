package fusecore

import (
	"context"
	"fmt"
	"fuse_file_system/log"
	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
	"github.com/sirupsen/logrus"
	"syscall"
)

type Dir struct {
	fs.Inode
	Attr  fuse.Attr
	child []fs.Inode
	name  string
}

func (d *Dir) Lookup(ctx context.Context, name string, out *fuse.EntryOut) (*fs.Inode, syscall.Errno) {
	logrus.Infoln("dir: lookup ", name, out)
	// todo
	inode := d.GetChild(name)
	if inode != nil {
		return inode, syscall.F_OK
	}
	return d.NewInode(ctx, &File{name: fmt.Sprintf("%s/%s", d.name, name)}, fs.StableAttr{Mode: fuse.S_IFREG}), syscall.F_OK
}

func (d *Dir) Mknod(ctx context.Context, name string, mode uint32, dev uint32, out *fuse.EntryOut) (*fs.Inode, syscall.Errno) {
	log.Logger.Info("dir mknod: ", name, mode, dev)
	newNode := d.NewInode(ctx, &File{}, fs.StableAttr{Mode: fuse.S_IFDIR})
	success := d.AddChild(name, newNode, false)
	if success {
		return newNode, syscall.ENOENT
	}
	return newNode, syscall.F_OK
}

func (d *Dir) Getattr(ctx context.Context, fh fs.FileHandle, out *fuse.AttrOut) syscall.Errno {
	logrus.Infoln("dir getattr: ", d.name)


	return syscall.F_OK
}
