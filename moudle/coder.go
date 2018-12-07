package moudle

import (
	"../def"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"hash/adler32"
	"io"
)

const headerSize = 8
const checksumSize = 8
const intSize = 8

type Decoder struct {
	mlen int
}

func NewDecoder() *Decoder {
	return &Decoder{
		mlen: 1024 * 1024 * 1,
	}
}

func NewDecoderWithSize(mlen int) *Decoder {
	return &Decoder{
		mlen: mlen,
	}
}

func (decoder *Decoder) ReadInt(reader io.Reader) (int, error) {
	_, err := decoder.readHeader(reader)
	if err != nil {
		return -1, err
	}

	valueBuf := make([]byte, intSize+checksumSize)

	_, err = reader.Read(valueBuf)

	if err != nil {
		return -1, err
	}

	checksumBuf := valueBuf[intSize:]
	valueBuf = valueBuf[:intSize]

	err = decoder.check(valueBuf, checksumBuf)

	if err != nil {
		return -1, err
	}

	return int(binary.LittleEndian.Uint32(valueBuf)), nil
}

func (decoder *Decoder) check(data []byte, checksumBuf []byte) error {

	reciveChecksum := adler32.Checksum(data)

	sendChecksum := binary.LittleEndian.Uint32(checksumBuf)

	if sendChecksum != reciveChecksum {
		return fmt.Errorf("checksum error")
	}

	return nil
}

func (decoder *Decoder) ReadStruct(reader io.Reader, v interface{}) error {
	length, err := decoder.readHeader(reader)
	if err != nil {
		return err
	}
	if length <= 0 || length > decoder.mlen {
		return fmt.Errorf("packet length error")
	}

	buf := make([]byte, length+checksumSize)

	_, err = reader.Read(buf)

	checksumBuf := buf[length:]

	buf = buf[:length]

	if err != nil {
		return err
	}

	err = decoder.check(buf, checksumBuf)

	if err != nil {
		return err
	}

	err = json.Unmarshal(buf, &v)

	if err != nil {
		return err
	}
	return nil
}

func (decoder *Decoder) readHeader(reader io.Reader) (int, error) {
	buf := make([]byte, headerSize)
	_, err := reader.Read(buf)
	if err != nil {
		return -1, err
	}

	length := binary.LittleEndian.Uint32(buf)

	return int(length), nil
}

type EnCoder struct {
	buf *bytes.Buffer
}

func NewEnCoder() *EnCoder {
	return &EnCoder{
		buf: bytes.NewBuffer(make([]byte, 0)),
	}
}

func (encoder *EnCoder) AppendInt(data int) {
	buf := bytes.NewBuffer(make([]byte, 0))

	binary.Write(buf, binary.LittleEndian, uint64(headerSize))

	binary.Write(buf, binary.LittleEndian, uint64(data))

	checksum := adler32.Checksum(buf.Bytes()[headerSize:])

	binary.Write(buf, binary.LittleEndian, uint64(checksum))

	encoder.buf.Write(buf.Bytes())
}

func (encoder *EnCoder) SendInt(writer io.Writer, data int) error {
	encoder.AppendInt(data)
	return encoder.Send(writer)
}

func (encoder *EnCoder) SendStruct(writer io.Writer, v def.SocketInterface) error {
	err := encoder.AppendStruct(v)
	if err != nil {
		return err
	}
	return encoder.Send(writer)
}

func (encoder *EnCoder) AppendStruct(v def.SocketInterface) error {
	data, err := v.StructToBytes()
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer(make([]byte, 0))

	binary.Write(buf, binary.LittleEndian, uint64(len(data)))

	binary.Write(buf, binary.LittleEndian, data)

	checksum := adler32.Checksum(data)

	binary.Write(buf, binary.LittleEndian, uint64(checksum))

	encoder.buf.Write(buf.Bytes())
	return nil
}

func (encoder *EnCoder) Send(writer io.Writer) error {
	_, err := writer.Write(encoder.buf.Bytes())
	if err != nil {
		return err
	}
	encoder.buf.Reset()
	return nil
}

type DEcoder struct {
	*Decoder
	*EnCoder
}

func NewDECoder() *DEcoder {
	return &DEcoder{
		NewDecoder(),
		NewEnCoder(),
	}
}

func NewDECoderWithSize(max_len int) *DEcoder {
	return &DEcoder{
		NewDecoderWithSize(max_len),
		NewEnCoder(),
	}
}
