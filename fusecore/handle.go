package fusecore

import (
	"context"
	"fuse_file_system/log"
	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
	"sync"
	"syscall"
)


type Handle struct {
	fd int
	name string
	mu sync.Mutex
}

func NewHandle(fd int,name string)*Handle{
	return &Handle{fd:fd,name:name}
}

func (fh *Handle) Flush(ctx context.Context) syscall.Errno {
	log.Logger.Infoln("handle flush: ",fh.name)
	fh.mu.Lock()
	defer fh.mu.Unlock()
	newFd, err := syscall.Dup(fh.fd)

	if err != nil {
		return fs.ToErrno(err)
	}
	err = syscall.Close(newFd)
	return fs.ToErrno(err)

}

func (fh *Handle) Release(ctx context.Context) syscall.Errno {
	log.Logger.Infoln("handle release: ", fh.name)
	fh.mu.Lock()
	defer fh.mu.Unlock()
	if fh.fd != -1 {
		err := syscall.Close(fh.fd)
		fh.fd = -1
		return fs.ToErrno(err)
	}
	return syscall.EBADF
}

func (fh *Handle) Write(ctx context.Context, data []byte, off int64) (written uint32, errno syscall.Errno) {
	log.Logger.Info("handle write: ", fh.name, len(data), off)

	return uint32(len(data)), syscall.F_OK
}

func (fh *Handle) Read(ctx context.Context, data []byte, off int64) (fuse.ReadResult, syscall.Errno) {
	log.Logger.Infoln("handle read: ",fh.name,len(data),off)
	fh.mu.Lock()
	defer fh.mu.Unlock()
	r := fuse.ReadResultFd(uintptr(fh.fd), off, len(data))
	return r, fs.OK
	//return &ReadResult{data: []byte("hello ..")}, fs.OK
}
