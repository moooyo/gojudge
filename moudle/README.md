#### 协议层模块

本模块包括:
1. socket: 带有buffer的net.Conn
2. decoder: 实现gojudge通信协议中的编码
3. encoder：实现gojudge通信协议的解码

### Socket
	socket完全实现了net.Conn接口，并在基础上添加了如下方法:
	1. func (socket *Socket)Flush() error 冲刷socket缓冲区
	2. func Dial(addr string) (*socket, error) 发起tcp连接
	3. func SocketFromConn(conn net.Conn) ×Socket 从conn构造Socket

```go
package moudle

import (
	"../def"
	"log"
	"net"
	"reflect"
	"testing"
)

func server(listen net.Listener) {
	for {
		conn, err := listen.Accept()
		if err != nil {
			continue
		}

		socket := SocketFromConn(conn)

		coder := NewDECoderWithSize(1024 * 1024 * 2)

		for i := 0; i < 1000; i++ {
			var submit def.Submit
			err = coder.ReadStruct(socket, &submit)
			if err != nil {
				log.Fatal(err)
			}

			coder.AppendStruct(&submit)

			coder.Send(socket)
		}

		socket.Close()
	}
}

func TestSocket(t *testing.T) {

	listen, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		t.Fatal(err)
	}

	go server(listen)

	socket, err := Dial("127.0.0.1:8080")
	if err != nil {
		t.Fatal(err)
	}

	coder := NewDECoderWithSize(2 * 1024 * 1024)

	defer socket.Close()

	for i := 0; i < 1000; i++ {
		submit := def.Submit{
			SubmitID:   i,
			ProblemID:  i,
			CodeSource: make([]byte, 1024*5),
			Language:   i,
		}
		coder.AppendStruct(&submit)

		coder.Send(socket)

		var result def.Submit

		err = coder.ReadStruct(socket, &result)

		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(result, submit) {
			t.Fail()
		}
	}
}
```


### Coder
```go
	//gojudge的消息格式
	type msg struct {
		len uint64
		data []byte
		checksum uint64
	}
```
#### Encoder
	动态添加数据，将数据编码成msg格式的[]byte存至内部buffer中，在需要写入时，调用send方法写入给定的writer

```go
type Encoder struct {
	buf *bytes.Buffer
}

//构造一个编码器
func NewEnCoder() *EnCoder {}

//向编码器buffer中追加一个整数
func (encoder *EnCoder) AppendInt(data int) {}

//向编码器中追加一个整数，并将编码器中所有数据写入writer
func (encoder *EnCoder) SendInt(writer io.Writer, data int) error {}

//向编码器中追加一个实现了SocketInterface接口的结构体，并将编码器内所有数据写入writer
func (encoder *EnCoder) SendStruct(writer io.Writer, v def.SocketInterface) error {}

//向编码器buffer中追加一个接口体
func (encoder *EnCoder) AppendStruct(v def.SocketInterface) error {}

//将编码器中所有数据写入writer
func (encoder *EnCoder) Send(writer io.Writer) error {}
```

#### Decoder

	从一个reader中读取长度不超过mlen的消息，并解码成特定的数据格式
```go
type Decoder struct {
	mlen int
}

// 构造能解不超过1m的消息解码器
func NewDecoder() *Decoder {}

// 构造能解不超过mlen的消息解码器
func NewDecoderWithSize(mlen int) *Decoder {}

//从reader中读取一个int
func (decoder *Decoder) ReadInt(reader io.Reader) (int, error) {}

//从reader中读取一个struct
func (decoder *Decoder) ReadStruct(reader io.Reader, v interface{}) error {}

```

```go
package moudle

import (
	"../def"
	"os"
	"reflect"
	"testing"
)

func TestCode(t *testing.T) {
	file, err := os.Create("test_code")
	if err != nil {
		t.Fatal(err)
	}
	coder := NewDECoderWithSize(1024 * 6)

	for i := 0; i < 10; i++ {
		submit := def.Submit{
			SubmitID:   i,
			ProblemID:  i,
			CodeSource: make([]byte, 4096),
			Language:   i,
		}
		err = coder.AppendStruct(&submit)
		coder.AppendInt(i)
		if err != nil {
			t.Fatal(err)
			file.Close()
			os.Remove(file.Name())
		}
	}
	err = coder.Send(file)
	if err != nil {
		file.Close()
		os.Remove(file.Name())
		t.Fatal(err)
	}

	file.Close()

	file, err = os.Open("test_code")

	defer file.Close()
	defer os.Remove(file.Name())
	for i := 0; i < 10; i++ {
		submit := def.Submit{
			SubmitID:   i,
			ProblemID:  i,
			CodeSource: make([]byte, 4096),
			Language:   i,
		}
		var result def.Submit

		err = coder.ReadStruct(file, &result)
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(result, submit) {
			t.Fail()
		}
		value, err := coder.ReadInt(file)
		if err != nil {
			t.Fatal(err)
		}
		if value != i {
			t.Fail()
		}
	}
}
```
