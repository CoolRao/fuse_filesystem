package dao

import "fuse_file_system/db"

type DirectDao struct {
	dao db.IDao
}

func (d *DirectDao) GetDirect() (string, error) {
	panic("implement me")
}

func (d *DirectDao) InsertDirect() (interface{}, error) {
	panic("implement me")
}

func (d *DirectDao) DirectList() ([]string, error) {
	panic("implement me")
}

func NewDirectDao() IDirectDao {
	return &DirectDao{dao: db.NewDao(db.DirectTable)}
}
