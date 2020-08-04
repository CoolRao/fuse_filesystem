package fusecore

import (
	"context"
	"fmt"
	"fuse_file_system/log"
	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
	"io/ioutil"
	"os"
	"syscall"
)

type Dir struct {
	fs.Inode
	Attr  fuse.Attr
	child []fs.Inode
	name  string
}

func (d *Dir) OnAdd(ctx context.Context) {
	path := fmt.Sprintf("%s/%s", copyDir, d.path())
	fileInfos, err := ioutil.ReadDir(path)
	if err != nil {
		log.Logger.Errorf("dir onAdd: %s", err.Error())
		return
	}
	for i := 0; i < len(fileInfos); i++ {
		info := fileInfos[i]
		if info.IsDir() {
			log.Logger.Debugf("dir OnAdd: %s ",info.Name())
			d.NewInode(ctx, &Dir{}, fs.StableAttr{Mode: fuse.S_IFDIR})
		} else {
			log.Logger.Debugf("dir OnAdd: %s ",info.Name())
			d.NewInode(ctx, &File{}, fs.StableAttr{Mode: fuse.S_IFREG})
		}
	}
}

func (d *Dir) Getattr(ctx context.Context, fh fs.FileHandle, out *fuse.AttrOut) syscall.Errno {
	log.Logger.Debugf("dir getattr:  %s ", d.name)
	out.Mode = fuse.S_IFDIR
	out.Size = 1024
	path := d.path()
	stat, err := os.Stat(fmt.Sprintf("%s/%s", copyDir, path))
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
	newNode := d.NewInode(ctx, &File{}, fs.StableAttr{Mode: fuse.S_IFDIR})
	return newNode, syscall.F_OK
}

func (d *Dir) Mkdir(ctx context.Context, name string, mode uint32, out *fuse.EntryOut) (*fs.Inode, syscall.Errno) {
	log.Logger.Debugf("dir mkdir:  %s %d %v", name, mode, out)
	newInode := d.NewInode(ctx, &Dir{}, fs.StableAttr{Mode: fuse.S_IFDIR})
	// todo
	err := os.Mkdir(fmt.Sprintf("%s/%s", copyDir, name), os.ModeDir)
	if err != nil {
		return nil, syscall.EIO
	}
	return newInode, syscall.F_OK
}
