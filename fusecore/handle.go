package fusecore

import (
	"context"
	"fuse_file_system/log"
	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
	"syscall"
)


type Handle struct {
	name string
}

func (fh *Handle) Flush(ctx context.Context) syscall.Errno {

	return syscall.F_OK
}

func (fh *Handle) Release(ctx context.Context) syscall.Errno {
	log.Logger.Infoln("handle release: ", fh.name)
	return syscall.F_OK
}

func (fh *Handle) Write(ctx context.Context, data []byte, off int64) (written uint32, errno syscall.Errno) {
	log.Logger.Info("handler write: ", fh.name, len(data), off)
	return uint32(len(data)), syscall.F_OK
}

func (fh *Handle) Read(ctx context.Context, data []byte, off int64) (fuse.ReadResult, syscall.Errno) {

	return &ReadResult{data: data}, fs.OK
}
