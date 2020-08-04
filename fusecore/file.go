package fusecore

import (
	"context"
	"fuse_file_system/log"
	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
	"path/filepath"
	"syscall"
	"time"
)

type File struct {
	fs.Inode
	name   string
	handle *Handle
}

func (f *File) path() string {
	path := f.Path(f.Root())
	return filepath.Join(f.root().rootPath, path)
}

func (n *File) root() *Root {
	return n.Root().Operations().(*Root)
}

func (f *File) Getattr(ctx context.Context, fh fs.FileHandle, out *fuse.AttrOut) syscall.Errno {
	log.Logger.Infoln("file getattr: ", f.name)
	out.Mode = fuse.S_IFREG
	out.Size = 1024
	out.Atime = uint64(time.Now().Unix())
	return syscall.F_OK
}

func (f *File) Write(ctx context.Context, data []byte, off int64) (written uint32, errno syscall.Errno) {
	log.Logger.Infoln("file write: ")
	return 1024, syscall.F_OK
}

func (f *File) Open(ctx context.Context, flags uint32) (fh fs.FileHandle, fuseFlags uint32, errno syscall.Errno) {
	lf := NewHandle(0, f.name)
	return lf, fuse.FOPEN_KEEP_CACHE, 0
}

func (f *File) Setattr(ctx context.Context, in *fuse.SetAttrIn, out *fuse.AttrOut) syscall.Errno {
	log.Logger.Infoln("file setattr: ", in, out)
	return 0
}
