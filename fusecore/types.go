package fusecore

import (
	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
	"syscall"
)

type ReadResult struct {
	data []byte
}

func (r *ReadResult) Size() int {
	return len(r.data)
}

func (r *ReadResult) Done() {

}
func (r *ReadResult) Bytes(buf []byte) ([]byte, fuse.Status) {
	return r.data, fuse.OK
}

type Direct struct {
	entries []fuse.DirEntry
}

func (d *Direct) HasNext() bool {
	return len(d.entries) > 0
}

func (d *Direct) Next() (fuse.DirEntry, syscall.Errno) {
	e := d.entries[0]
	d.entries = d.entries[1:]
	return e, syscall.F_OK
}

func (d *Direct) Close() {

}

func NewDirect(param map[string]*fs.Inode) *Direct {
	var res []fuse.DirEntry
	for k, v := range param {
		res = append(res, fuse.DirEntry{
			Mode: v.Mode(),
			Name: k,
			Ino:  v.StableAttr().Ino,
		})
	}
	return &Direct{res}
}


