package db

import (
	"fuse_file_system/model"
	"log"
)

type Dao struct {
}

func (d *Dao) Delete(param model.Param) (interface{}, error) {
	// todo
	stmt, err := db.Prepare("delete from userinfo where uid=?")
	if err!=nil{
		return nil,err
	}
	res, err := stmt.Exec(param)
	if err!=nil{
		return nil,err
	}
	affect, err := res.RowsAffected()
	if err!=nil{
		return nil,err
	}
	return affect,nil
}

func (d *Dao) Insert(param model.Param) (interface{}, error) {
	// todo
	stmt, err := GetDb().Prepare("insert into foo(id, name) values(?, ?)")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	_, err = stmt.Exec("id0001","value001")
	if err != nil {
		return nil,err
	}
	return nil,nil
}

func (d *Dao) FindOne(param model.Param) (string, error) {
	// todo
	stmt, err := GetDb().Prepare("select name from foo where id = ?")
	if err != nil {
		return "", err
	}
	defer stmt.Close()
	var name string
	err = stmt.QueryRow(param).Scan(&name)
	if err != nil {
		log.Fatal(err)
	}
	return name, nil
}

func NewDao() IDao {
	return &Dao{}
}


func findSql(param model.Param) string{
	panic(param)
}
