package fusecore

import (
	"context"
	"github.com/hanwen/go-fuse/v2/fs"
	"golang.org/x/sys/unix"
	"os"
	"path/filepath"
	"syscall"

	"github.com/hanwen/go-fuse/v2/fuse"
)

type loopbackRoot struct {
	loopbackNode

	rootPath string
	rootDev  uint64
}

type loopbackNode struct {
	fs.Inode
}

var _ = (fs.NodeStatfser)((*loopbackNode)(nil))
var _ = (fs.NodeStatfser)((*loopbackNode)(nil))
var _ = (fs.NodeGetattrer)((*loopbackNode)(nil))
var _ = (fs.NodeReadlinker)((*loopbackNode)(nil))
var _ = (fs.NodeOpener)((*loopbackNode)(nil))
var _ = (fs.NodeLookuper)((*loopbackNode)(nil))
var _ = (fs.NodeOpendirer)((*loopbackNode)(nil))
var _ = (fs.NodeReaddirer)((*loopbackNode)(nil))
var _ = (fs.NodeMkdirer)((*loopbackNode)(nil))
var _ = (fs.NodeMknoder)((*loopbackNode)(nil))
var _ = (fs.NodeLinker)((*loopbackNode)(nil))
var _ = (fs.NodeSymlinker)((*loopbackNode)(nil))
var _ = (fs.NodeUnlinker)((*loopbackNode)(nil))
var _ = (fs.NodeRmdirer)((*loopbackNode)(nil))
var _ = (fs.NodeRenamer)((*loopbackNode)(nil))

func (n *loopbackNode) Statfs(ctx context.Context, out *fuse.StatfsOut) syscall.Errno {
	s := syscall.Statfs_t{}
	err := syscall.Statfs(n.path(), &s)
	if err != nil {
		return fs.ToErrno(err)
	}
	out.FromStatfsT(&s)
	return fs.OK
}

func (r *loopbackRoot) Getattr(ctx context.Context, f fs.FileHandle, out *fuse.AttrOut) syscall.Errno {
	st := syscall.Stat_t{}
	err := syscall.Stat(r.rootPath, &st)
	if err != nil {
		return fs.ToErrno(err)
	}
	out.FromStat(&st)
	return fs.OK
}

func (n *loopbackNode) root() *loopbackRoot {
	return n.Root().Operations().(*loopbackRoot)
}

func (n *loopbackNode) path() string {
	path := n.Path(n.Root())
	return filepath.Join(n.root().rootPath, path)
}

func (n *loopbackNode) Lookup(ctx context.Context, name string, out *fuse.EntryOut) (*fs.Inode, syscall.Errno) {
	p := filepath.Join(n.path(), name)

	st := syscall.Stat_t{}
	err := syscall.Lstat(p, &st)
	if err != nil {
		return nil, fs.ToErrno(err)
	}

	out.Attr.FromStat(&st)
	node := &loopbackNode{}
	ch := n.NewInode(ctx, node, n.root().idFromStat(&st))
	return ch, 0
}

// preserveOwner sets uid and gid of `rootPath` according to the caller information
// in `ctx`.
func (n *loopbackNode) preserveOwner(ctx context.Context, path string) error {
	if os.Getuid() != 0 {
		return nil
	}
	caller, ok := fuse.FromContext(ctx)
	if !ok {
		return nil
	}
	return syscall.Lchown(path, int(caller.Uid), int(caller.Gid))
}

func (n *loopbackNode) Mknod(ctx context.Context, name string, mode, rdev uint32, out *fuse.EntryOut) (*fs.Inode, syscall.Errno) {
	p := filepath.Join(n.path(), name)
	err := syscall.Mknod(p, mode, int(rdev))
	if err != nil {
		return nil, fs.ToErrno(err)
	}
	n.preserveOwner(ctx, p)
	st := syscall.Stat_t{}
	if err := syscall.Lstat(p, &st); err != nil {
		syscall.Rmdir(p)
		return nil, fs.ToErrno(err)
	}

	out.Attr.FromStat(&st)

	node := &loopbackNode{}
	ch := n.NewInode(ctx, node, n.root().idFromStat(&st))

	return ch, 0
}

func (n *loopbackNode) Mkdir(ctx context.Context, name string, mode uint32, out *fuse.EntryOut) (*fs.Inode, syscall.Errno) {
	p := filepath.Join(n.path(), name)
	err := os.Mkdir(p, os.FileMode(mode))
	if err != nil {
		return nil, fs.ToErrno(err)
	}
	n.preserveOwner(ctx, p)
	st := syscall.Stat_t{}
	if err := syscall.Lstat(p, &st); err != nil {
		syscall.Rmdir(p)
		return nil, fs.ToErrno(err)
	}

	out.Attr.FromStat(&st)

	node := &loopbackNode{}
	ch := n.NewInode(ctx, node, n.root().idFromStat(&st))

	return ch, 0
}

func (n *loopbackNode) Rmdir(ctx context.Context, name string) syscall.Errno {
	p := filepath.Join(n.path(), name)
	err := syscall.Rmdir(p)
	return fs.ToErrno(err)
}

func (n *loopbackNode) Unlink(ctx context.Context, name string) syscall.Errno {
	p := filepath.Join(n.path(), name)
	err := syscall.Unlink(p)
	return fs.ToErrno(err)
}

func toLoopbackNode(op fs.InodeEmbedder) *loopbackNode {
	if r, ok := op.(*loopbackRoot); ok {
		return &r.loopbackNode
	}
	return op.(*loopbackNode)
}

func (n *loopbackNode) Rename(ctx context.Context, name string, newParent fs.InodeEmbedder, newName string, flags uint32) syscall.Errno {
	newParentLoopback := toLoopbackNode(newParent)
	if flags&fs.RENAME_EXCHANGE != 0 {
		return n.renameExchange(name, newParentLoopback, newName)
	}

	p1 := filepath.Join(n.path(), name)
	p2 := filepath.Join(newParentLoopback.path(), newName)

	err := syscall.Rename(p1, p2)
	return fs.ToErrno(err)
}

func (n *loopbackNode) renameExchange(name string, newparent *loopbackNode, newName string) syscall.Errno {
	fd1, err := syscall.Open(n.path(), syscall.O_DIRECTORY, 0)
	if err != nil {
		return fs.ToErrno(err)
	}
	defer syscall.Close(fd1)
	fd2, err := syscall.Open(newparent.path(), syscall.O_DIRECTORY, 0)
	defer syscall.Close(fd2)
	if err != nil {
		return fs.ToErrno(err)
	}

	var st syscall.Stat_t
	if err := syscall.Fstat(fd1, &st); err != nil {
		return fs.ToErrno(err)
	}

	inode := &n.Inode
	if inode.Root() != inode && inode.StableAttr().Ino != n.root().idFromStat(&st).Ino {
		return syscall.EBUSY
	}
	if err := syscall.Fstat(fd2, &st); err != nil {
		return fs.ToErrno(err)
	}

	newinode := &newparent.Inode
	if newinode.Root() != newinode && newinode.StableAttr().Ino != n.root().idFromStat(&st).Ino {
		return syscall.EBUSY
	}

	return fs.ToErrno(unix.Renameat2(fd1, name, fd2, newName, unix.RENAME_EXCHANGE))
}


func (r *loopbackRoot) idFromStat(st *syscall.Stat_t) fs.StableAttr {
	swapped := (uint64(st.Dev) << 32) | (uint64(st.Dev) >> 32)
	swappedRootDev := (r.rootDev << 32) | (r.rootDev >> 32)
	return fs.StableAttr{
		Mode: uint32(st.Mode),
		Gen:  1,
		Ino: (swapped ^ swappedRootDev) ^ st.Ino,
	}
}

var _ = (fs.NodeCreater)((*loopbackNode)(nil))

func (n *loopbackNode) Create(ctx context.Context, name string, flags uint32, mode uint32, out *fuse.EntryOut) (inode *fs.Inode, fh fs.FileHandle, fuseFlags uint32, errno syscall.Errno) {
	p := filepath.Join(n.path(), name)
	flags = flags &^ syscall.O_APPEND
	fd, err := syscall.Open(p, int(flags)|os.O_CREATE, mode)
	if err != nil {
		return nil, nil, 0, fs.ToErrno(err)
	}
	n.preserveOwner(ctx, p)
	st := syscall.Stat_t{}
	if err := syscall.Fstat(fd, &st); err != nil {
		syscall.Close(fd)
		return nil, nil, 0, fs.ToErrno(err)
	}

	node := &loopbackNode{}
	ch := n.NewInode(ctx, node, n.root().idFromStat(&st))
	lf := fs.NewLoopbackFile(fd)

	out.FromStat(&st)
	return ch, lf, 0, 0
}

func (n *loopbackNode) Symlink(ctx context.Context, target, name string, out *fuse.EntryOut) (*fs.Inode, syscall.Errno) {
	p := filepath.Join(n.path(), name)
	err := syscall.Symlink(target, p)
	if err != nil {
		return nil, fs.ToErrno(err)
	}
	n.preserveOwner(ctx, p)
	st := syscall.Stat_t{}
	if err := syscall.Lstat(p, &st); err != nil {
		syscall.Unlink(p)
		return nil, fs.ToErrno(err)
	}
	node := &loopbackNode{}
	ch := n.NewInode(ctx, node, n.root().idFromStat(&st))

	out.Attr.FromStat(&st)
	return ch, 0
}

func (n *loopbackNode) Link(ctx context.Context, target fs.InodeEmbedder, name string, out *fuse.EntryOut) (*fs.Inode, syscall.Errno) {

	p := filepath.Join(n.path(), name)
	targetNode := toLoopbackNode(target)
	err := syscall.Link(targetNode.path(), p)
	if err != nil {
		return nil, fs.ToErrno(err)
	}
	st := syscall.Stat_t{}
	if err := syscall.Lstat(p, &st); err != nil {
		syscall.Unlink(p)
		return nil, fs.ToErrno(err)
	}
	node := &loopbackNode{}
	ch := n.NewInode(ctx, node, n.root().idFromStat(&st))

	out.Attr.FromStat(&st)
	return ch, 0
}

func (n *loopbackNode) Readlink(ctx context.Context) ([]byte, syscall.Errno) {
	p := n.path()

	for l := 256; ; l *= 2 {
		buf := make([]byte, l)
		sz, err := syscall.Readlink(p, buf)
		if err != nil {
			return nil, fs.ToErrno(err)
		}

		if sz < len(buf) {
			return buf[:sz], 0
		}
	}
}

func (n *loopbackNode) Open(ctx context.Context, flags uint32) (fh fs.FileHandle, fuseFlags uint32, errno syscall.Errno) {
	flags = flags &^ syscall.O_APPEND
	p := n.path()
	f, err := syscall.Open(p, int(flags), 0)
	if err != nil {
		return nil, 0, fs.ToErrno(err)
	}
	lf := fs.NewLoopbackFile(f)
	return lf, 0, 0
}

func (n *loopbackNode) Opendir(ctx context.Context) syscall.Errno {
	fd, err := syscall.Open(n.path(), syscall.O_DIRECTORY, 0755)
	if err != nil {
		return fs.ToErrno(err)
	}
	syscall.Close(fd)
	return fs.OK
}

func (n *loopbackNode) Readdir(ctx context.Context) (fs.DirStream, syscall.Errno) {
	return fs.NewLoopbackDirStream(n.path())
}

func (n *loopbackNode) Getattr(ctx context.Context, f fs.FileHandle, out *fuse.AttrOut) syscall.Errno {
	if f != nil {
		return f.(fs.FileGetattrer).Getattr(ctx, out)
	}
	p := n.path()

	var err error
	st := syscall.Stat_t{}
	err = syscall.Lstat(p, &st)
	if err != nil {
		return fs.ToErrno(err)
	}
	out.FromStat(&st)
	return fs.OK
}

var _ = (fs.NodeSetattrer)((*loopbackNode)(nil))

func (n *loopbackNode) Setattr(ctx context.Context, f fs.FileHandle, in *fuse.SetAttrIn, out *fuse.AttrOut) syscall.Errno {
	p := n.path()
	fsa, ok := f.(fs.FileSetattrer)
	if ok && fsa != nil {
		fsa.Setattr(ctx, in, out)
	} else {
		if m, ok := in.GetMode(); ok {
			if err := syscall.Chmod(p, m); err != nil {
				return fs.ToErrno(err)
			}
		}

		uid, uok := in.GetUID()
		gid, gok := in.GetGID()
		if uok || gok {
			suid := -1
			sgid := -1
			if uok {
				suid = int(uid)
			}
			if gok {
				sgid = int(gid)
			}
			if err := syscall.Chown(p, suid, sgid); err != nil {
				return fs.ToErrno(err)
			}
		}

		mtime, mok := in.GetMTime()
		atime, aok := in.GetATime()

		if mok || aok {

			ap := &atime
			mp := &mtime
			if !aok {
				ap = nil
			}
			if !mok {
				mp = nil
			}
			var ts [2]syscall.Timespec
			ts[0] = fuse.UtimeToTimespec(ap)
			ts[1] = fuse.UtimeToTimespec(mp)

			if err := syscall.UtimesNano(p, ts[:]); err != nil {
				return fs.ToErrno(err)
			}
		}

		if sz, ok := in.GetSize(); ok {
			if err := syscall.Truncate(p, int64(sz)); err != nil {
				return fs.ToErrno(err)
			}
		}
	}

	fga, ok := f.(fs.FileGetattrer)
	if ok && fga != nil {
		fga.Getattr(ctx, out)
	} else {
		st := syscall.Stat_t{}
		err := syscall.Lstat(p, &st)
		if err != nil {
			return fs.ToErrno(err)
		}
		out.FromStat(&st)
	}
	return fs.OK
}

// NewLoopbackRoot returns a root node for a loopback file system whose
// root is at the given root. This node implements all NodeXxxxer
// operations available.
func NewLoopbackRoot(root string) (fs.InodeEmbedder, error) {
	var st syscall.Stat_t
	err := syscall.Stat(root, &st)
	if err != nil {
		return nil, err
	}

	n := &loopbackRoot{
		rootPath: root,
		rootDev:  uint64(st.Dev),
	}
	return n, nil
}
