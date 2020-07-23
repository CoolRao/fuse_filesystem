package db

import (
	"bytes"
	"fmt"
	"fuse_file_system/model"
	"log"
)

type Dao struct {
	tableName string
}

func (d *Dao) List(sql string, param []interface{}) (interface{}, error) {
	stmt, err := GetDb().Prepare(sql)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	res, err := stmt.Query(param)
	if err != nil {
		return nil, err
	}
	// todo
	fields, _ := res.Columns()
	fieldslen := len(fields)
	ms := []map[string]interface{}{}
	for res.Next() {
		m := map[string]interface{}{}
		scaninterfaces := make([]interface{}, 0, fieldslen)
		for _, field := range fields {
			//m[field] = 0
			scaninterfaces = append(scaninterfaces, m[field])
		}
		res.Scan(scaninterfaces...)
		ms = append(ms, m)
	}
	return ms, nil
}

func (d *Dao) getTable() string {
	return d.tableName
}

func NewDao(table string) IDao {
	return &Dao{tableName: table}
}

func (d *Dao) Delete(sql string, param []interface{}) (interface{}, error) {
	stmt, err := db.Prepare(sql)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	res, err := stmt.Exec(param)
	if err != nil {
		return nil, err
	}
	affect, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	return affect, nil
}

func (d *Dao) Insert(sql string, param []interface{}) (interface{}, error) {
	stmt, err := GetDb().Prepare(sql)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	_, err = stmt.Exec(param)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (d *Dao) FindOne(sql string, param []interface{}, values...interface{} ) error {
	stmt, err := GetDb().Prepare(sql)
	if err != nil {
		return err
	}
	defer stmt.Close()
	err = stmt.QueryRow(param).Scan(values)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func deleteSql(table string, param model.Param) (string, []interface{}) {
	var keys []string
	var values []interface{}
	for k, v := range param {
		keys = append(keys, k)
		values = append(values, v)
	}
	var keyBuffer bytes.Buffer
	var valueBuffer bytes.Buffer
	valueBuffer.WriteString("(")
	for i := 0; i < len(keys); i++ {
		if i == 0 {
			keyBuffer.WriteString(fmt.Sprintf("%s = ?", keys[i]))
		} else {
			keyBuffer.WriteString(fmt.Sprintf("and %s = ?", keys[i]))
		}
	}
	return fmt.Sprintf("delete from %s where %s", table, keyBuffer.String()), values
}

func querySql(table string, param model.Param) (string, []interface{}) {
	var keys []string
	var values []interface{}
	for k, v := range param {
		keys = append(keys, k)
		values = append(values, v)
	}
	var keyBuffer bytes.Buffer
	var valueBuffer bytes.Buffer
	valueBuffer.WriteString("(")
	for i := 0; i < len(keys); i++ {
		if i == 0 {
			keyBuffer.WriteString(fmt.Sprintf("%s = ?", keys[i]))
		} else {
			keyBuffer.WriteString(fmt.Sprintf("and %s = ?", keys[i]))
		}
	}
	return fmt.Sprintf("select * from %s where %s", table, keyBuffer.String()), values
}

func insertSql(table string, param model.Param) (string, []interface{}) {
	var keys []string
	var values []interface{}
	for k, v := range param {
		keys = append(keys, k)
		values = append(values, v)
	}
	var keyBuffer bytes.Buffer
	var valueBuffer bytes.Buffer

	keyBuffer.WriteString("(")
	valueBuffer.WriteString("(")
	for i := 0; i < len(keys); i++ {
		if i == len(keys)-1 {
			keyBuffer.WriteString(fmt.Sprintf("%s)", keys[i]))
			valueBuffer.WriteString(fmt.Sprintf("?)"))
		} else {
			keyBuffer.WriteString(fmt.Sprintf("%s,", keys[i]))
			valueBuffer.WriteString(fmt.Sprintf("?,"))
		}
	}
	return fmt.Sprintf("insert into %s%s values%s", table, keyBuffer.String(), valueBuffer.String()), values
}
