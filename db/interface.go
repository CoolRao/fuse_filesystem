package db

import "fuse_file_system/model"

type IDao interface {
	FindOne(param model.Param)(string,error)
	Insert(param model.Param)(interface{},error)
	Delete(param model.Param)(interface{},error)
}
