package fusecore

import (
	"context"
	"fuse_file_system/log"
	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
	"github.com/sirupsen/logrus"
	"syscall"
)

type File struct {
	fs.Inode
	name string
}

func (f *File) Getattr(ctx context.Context, fh fs.FileHandle, out *fuse.AttrOut) syscall.Errno {
	return syscall.F_OK
}

func (f *File) Open(ctx context.Context, flags uint32) (fh fs.FileHandle, fuseFlags uint32, errno syscall.Errno) {
	log.Logger.Infoln("file open: ", f.name)
	return &Handle{name: f.name}, fuse.FOPEN_KEEP_CACHE, fs.OK
}
func (f *File) Setattr(ctx context.Context, in *fuse.SetAttrIn, out *fuse.AttrOut) syscall.Errno {
	logrus.Infoln("file open: ", in, out)
	return 0
}

