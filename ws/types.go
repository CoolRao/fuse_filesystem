package ws

const (
	FileStateType = iota + 1
	FileReadType
	FileSyncAttr
)

type SendMsg struct {
	Type int
	Body []byte
}

type RecvMsg struct {
	Type int
	Body []byte
}

type Message struct {
	MsgType  int
	FileName string
	DirType  string
	MsgId    string
	Body     []byte
	Size     int64
	Off      int64
}
