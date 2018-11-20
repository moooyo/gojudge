package def

type SocketInterface interface {
	StructToBytes() (data []byte,err error)
}
