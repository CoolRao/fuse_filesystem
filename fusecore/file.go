package fusecore

import (
	"context"
	"fuse_file_system/log"
	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
	"path/filepath"
	"syscall"
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

func (f *File) root() *Root {
	return f.Root().Operations().(*Root)
}

func (f *File) Getattr(ctx context.Context, fh fs.FileHandle, out *fuse.AttrOut) syscall.Errno {
	log.Logger.Debugf("file getattr: %s ", f.name)
	//stat, err := os.Stat(getRemotePath(f.name))
	//if err != nil {
	//	return syscall.EIO
	//}
	stat, err := FileState(f.name)
	if err != nil {
		return syscall.EIO
	}
	out.Mode = uint32(stat.Mode)
	out.Size = uint64(stat.Size)
	out.Mtime = uint64(stat.ModTime.Unix())
	return syscall.F_OK
}

func (f *File) Open(ctx context.Context, flags uint32) (fh fs.FileHandle, fuseFlags uint32, errno syscall.Errno) {
	log.Logger.Debugf("file Open:  %s", f.name)
	flags = flags &^ syscall.O_APPEND
	lf := NewHandle(0, f.name)
	return lf, flags, 0
}

func (f *File) Setattr(ctx context.Context, in *fuse.SetAttrIn, out *fuse.AttrOut) syscall.Errno {
	log.Logger.Debugf("file setattr:  %v  %v ", in, out)
	return 0
}
