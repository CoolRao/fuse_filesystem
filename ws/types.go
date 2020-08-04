package ws

type SendMsg struct {
	Type int
	Body []byte
}

type RecvMsg struct {
	Type int
	Body []byte
}
