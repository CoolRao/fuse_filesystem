package db

type IDao interface {
	FindOne(sql string, param []interface{}, values ...interface{}) error
	Insert(sql string, param []interface{}) (interface{}, error)
	Delete(sql string, param []interface{}) (interface{}, error)
	List(sql string, param []interface{}) (interface{}, error)
}
