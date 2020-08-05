package fusecore

import (
	"context"
	"fuse_file_system/log"
	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
	"os"
	"syscall"
	"time"
)

type Dir struct {
	fs.Inode
	Attr  fuse.Attr
	child []fs.Inode
	name  string
}

func (d *Dir) OnAdd(ctx context.Context) {
	log.Logger.Debugf("dir onAdd: %s ", d.name)
}

func (d *Dir) Getattr(ctx context.Context, fh fs.FileHandle, out *fuse.AttrOut) syscall.Errno {
	log.Logger.Debugf("dir getattr:  %s ", d.name)
	path := d.path()
	stat, err := os.Stat(getRemotePath(path))
	if err != nil {
		return syscall.EIO
	}
	out.Size = uint64(stat.Size())
	out.Mode = uint32(stat.Mode())
	out.Mtime = uint64(stat.ModTime().Unix())
	return syscall.F_OK
}

func (d *Dir) path() string {
	return d.Path(d.Root())
}

func (d *Dir) Create(ctx context.Context, name string, flags uint32, mode uint32, out *fuse.EntryOut) (node *fs.Inode, fh fs.FileHandle, fuseFlags uint32, errno syscall.Errno) {
	log.Logger.Debugf("dir Create: %s  %d %d", name, flags, mode)
	file, err := os.Create(getRemotePath(name))
	if err != nil {
		return nil, nil, 0, syscall.EIO
	}
	stat, err := file.Stat()
	if err!=nil{
		return nil,nil,0,syscall.EIO
	}
	out.Mode= uint32(stat.Mode())
	out.Size= uint64(stat.Size())
	out.Mtime= uint64(stat.ModTime().Unix())
	out.Ctime= uint64(time.Now().Unix())
	out.Atime= uint64(time.Now().Unix())
	newInode := d.NewInode(ctx, &File{name:name}, fs.StableAttr{Mode: fuse.S_IFREG})
	return newInode, Handle{}, 0, syscall.F_OK
}

func (d *Dir) Lookup(ctx context.Context, name string, out *fuse.EntryOut) (*fs.Inode, syscall.Errno) {
	log.Logger.Debug("dir lookup: ", name, out)
	inode := d.GetChild(name)
	if inode != nil {
		return inode, syscall.F_OK
	}
	return nil, syscall.ENOENT
}

func (d *Dir) Mknod(ctx context.Context, name string, mode uint32, dev uint32, out *fuse.EntryOut) (*fs.Inode, syscall.Errno) {
	log.Logger.Debugf("dir mknod:  %s %d %d", name, mode, dev)
	newNode := d.NewInode(ctx, &File{name:name}, fs.StableAttr{Mode: fuse.S_IFDIR})
	return newNode, syscall.F_OK
}

func (d *Dir) Mkdir(ctx context.Context, name string, mode uint32, out *fuse.EntryOut) (*fs.Inode, syscall.Errno) {
	log.Logger.Debugf("dir mkdir:  %s %d %v", name, mode, out)
	newInode := d.NewInode(ctx, &Dir{}, fs.StableAttr{Mode: fuse.S_IFDIR})
	err := os.Mkdir(getRemotePath(name), os.ModeDir)
	if err != nil {
		return nil, syscall.EIO
	}
	return newInode, syscall.F_OK
}
