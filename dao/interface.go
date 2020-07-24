package dao

type IDirectDao interface {
	GetDirect() (string, error)
	InsertDirect() (interface{}, error)
	DirectList() ([]string, error)
}
