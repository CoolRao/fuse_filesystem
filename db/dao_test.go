package db

import (
	"fmt"
	"fuse_file_system/model"
	"testing"
)

func TestDao_InsertSql(t *testing.T) {
	param:=model.Param{}
	param["rao"]="xjrw"
	sql := insertSql("test", param)
	fmt.Println(sql)
	s, i := querySql("test", param)
	fmt.Println(s,i)
	s2, i2 := deleteSql("test", param)
	fmt.Println(s2,i2)
}
