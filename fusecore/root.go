package fusecore

import (
	"context"
	"fuse_file_system/log"
	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
	"syscall"
)

type Root struct {
	Dir
	rootDev  uint64
	rootPath string
	dest     string
}

func NewRoot(path string, dest string) *Root {
	return &Root{
		rootPath: path,
		dest:     dest,
	}
}

func (r *Root) Getattr(ctx context.Context, f fs.FileHandle, out *fuse.AttrOut) syscall.Errno {
	log.Logger.Infoln("root getattr: ", out)
	out.Mode = fuse.S_IFDIR
	return syscall.F_OK

}

func (r *Root) Statfs(ctx context.Context, out *fuse.StatfsOut) syscall.Errno {

	return syscall.F_OK
}
