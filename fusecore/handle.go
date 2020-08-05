package fusecore

import (
	"context"
	"fuse_file_system/log"
	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
	"os"
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
	log.Logger.Debugf("handle flush:  %s ",fh.name)
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
	log.Logger.Debugf("handle release:  %s", fh.name)
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
	log.Logger.Debugf("handle write:  %s  %d  %d  %s ", fh.name, len(data), off,string(data))
	path := getRemotePath(fh.name)
	log.Logger.Warnln(path)

	file, err := os.OpenFile(path,os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err!=nil{
		return 0,syscall.EIO
	}
	_, err = file.Seek(off, 0)
	if err!=nil{
		return 0,syscall.EIO
	}
	defer file.Close()
	_, err = file.Write(data)
	if err!=nil{
		return 0,syscall.EIO
	}
	return uint32(len(data)), syscall.F_OK
}

func (fh *Handle) Read(ctx context.Context, data []byte, off int64) (fuse.ReadResult, syscall.Errno) {
	log.Logger.Debugf("handle read:  %s  %d  %d ",fh.name,len(data),off)
	fh.mu.Lock()
	defer fh.mu.Unlock()
	path := getRemotePath(fh.name)
	file, err := os.Open(path)
	if err!=nil{
		return nil,syscall.EIO
	}
	_, err = file.Seek(off, 0)
	if err!=nil{
		return nil,syscall.EIO
	}
	defer file.Close()

	bytes:=make([]byte,len(data))
	_, err = file.Read(bytes)
	if err!=nil{
		return nil,syscall.EIO
	}
	return &ReadResult{data: bytes}, syscall.F_OK
}
