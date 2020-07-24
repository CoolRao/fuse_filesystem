package fusecore

import (
	"context"
	"fuse_file_system/log"
	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
	"math"
	"syscall"
	"time"
)

type Root struct {
	fs.Inode
}

func NewRoot() *Root {
	return &Root{}
}

func (r *Root) Getattr(ctx context.Context, f fs.FileHandle, out *fuse.AttrOut) syscall.Errno {
	log.Logger.Infoln("root getattr: ", out)
	size := math.Pow(1024, 6)
	out.Mode = fuse.S_IFDIR
	out.Size = uint64(size)
	out.Atime = uint64(time.Now().Unix())
	return syscall.F_OK

}

func (r *Root) Lookup(ctx context.Context, name string, out *fuse.EntryOut) (*fs.Inode, syscall.Errno) {
	log.Logger.Infoln("root lookup: ", name)
	child := r.GetChild(name)
	if child == nil {
		return nil, syscall.ENOENT
	}
	return child, syscall.F_OK

}

// 这个是创建并打开一个文件
func (r *Root) Create(ctx context.Context, name string, flags uint32, mode uint32, out *fuse.EntryOut) (node *fs.Inode, fh fs.FileHandle, fuseFlags uint32, errno syscall.Errno) {
	log.Logger.Infoln("Root create: ", name, flags, mode, out)
	inode := r.NewInode(ctx, &Dir{}, fs.StableAttr{Mode: fuse.S_IFREG})
	r.AddChild(name, inode, false)
	return inode, nil, flags, syscall.F_OK

}

func (r *Root) Mkdir(ctx context.Context, name string, mode uint32, out *fuse.EntryOut) (*fs.Inode, syscall.Errno) {
	log.Logger.Infoln("Root mkdir :", name)
	out.Mode = mode
	inode := r.NewInode(ctx, &Dir{}, fs.StableAttr{Mode: fuse.S_IFDIR})
	return inode, syscall.F_OK
}

func (r *Root) Mknod(ctx context.Context, name string, mode uint32, dev uint32, out *fuse.EntryOut) (*fs.Inode, syscall.Errno) {
	log.Logger.Infoln("Root: mknod ", name)
	newNode := r.NewInode(ctx, &File{}, fs.StableAttr{Mode: fuse.S_IFDIR})
	success := r.AddChild(name, newNode, false)
	if success {
		return newNode, syscall.ENOENT
	}
	return newNode, syscall.F_OK
}

func (r *Root) Readdir(ctx context.Context) (fs.DirStream, syscall.Errno) {
	children := r.Children()
	log.Logger.Infoln("root readdir: ", children)
	return NewDirect(children), syscall.F_OK
}

// 定义该文件系统相关信息
func (r *Root) Statfs(ctx context.Context, out *fuse.StatfsOut) syscall.Errno {
	log.Logger.Infoln("Root statfs: ", )
	out.Bavail = 1024
	out.Bfree = 1024
	out.Blocks = 4098
	out.Bsize = 1024
	out.Ffree = 1024
	out.Files = 1024
	out.Frsize = 1024
	out.NameLen = 1024
	out.Padding = 1024
	out.Spare = [6]uint32{}
	return 0
}

func (r *Root) OnAdd(ctx context.Context) {

}
