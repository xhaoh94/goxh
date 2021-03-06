package codec

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"io/ioutil"

	"github.com/xhaoh94/goxh/app"
	"github.com/xhaoh94/goxh/consts"
)

type (
	ICodec interface {
		Encode(interface{}) ([]byte, error)
		Decode([]byte, interface{}) error
	}
)

var codec ICodec

func SetCodec(icodec ICodec) {
	codec = icodec
}
func Encode(msg interface{}) ([]byte, error) {
	if msg == nil {
		return nil, consts.CodecError
	}
	return codec.Encode(msg)
}
func Decode(bytes []byte, msg interface{}) error {
	if msg == nil {
		return consts.CodecError
	}
	return codec.Decode(bytes, msg)
}

//BytesToUint16 转uint16
func BytesToUint16(b []byte) uint16 {
	bytesBuffer := bytes.NewBuffer(b)
	var tmp uint16
	binary.Read(bytesBuffer, app.NetEndian, &tmp)
	return tmp
}

//BytesToint16 转int16
func BytesToint16(b []byte) int16 {
	bytesBuffer := bytes.NewBuffer(b)
	var tmp int16
	binary.Read(bytesBuffer, app.NetEndian, &tmp)
	return tmp
}

//BytesToUint32 转uint32
func BytesToUint32(b []byte) uint32 {
	bytesBuffer := bytes.NewBuffer(b)
	var tmp uint32
	binary.Read(bytesBuffer, app.NetEndian, &tmp)
	return tmp
}

//BytesToint32 转int32
func BytesToint32(b []byte) int32 {
	bytesBuffer := bytes.NewBuffer(b)
	var tmp int32
	binary.Read(bytesBuffer, app.NetEndian, &tmp)
	return tmp
}

//BytesToUint64 转uint64
func BytesToUint64(b []byte) uint64 {
	bytesBuffer := bytes.NewBuffer(b)
	var tmp uint64
	binary.Read(bytesBuffer, app.NetEndian, &tmp)
	return tmp
}

//BytesToint64 转int64
func BytesToint64(b []byte) int64 {
	bytesBuffer := bytes.NewBuffer(b)
	var tmp int64
	binary.Read(bytesBuffer, app.NetEndian, &tmp)
	return tmp
}

//Uint16ToBytes 转bytes
func Uint16ToBytes(n uint16) []byte {
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, app.NetEndian, &n)
	return bytesBuffer.Bytes()
}

//Int16ToBytes 转bytes
func Int16ToBytes(n int16) []byte {
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, app.NetEndian, &n)
	return bytesBuffer.Bytes()
}

//Uint32ToBytes 转bytes
func Uint32ToBytes(n uint32) []byte {
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, app.NetEndian, &n)
	return bytesBuffer.Bytes()
}

//Int32ToBytes 转bytes
func Int32ToBytes(n int32) []byte {
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, app.NetEndian, &n)
	return bytesBuffer.Bytes()
}

//Uint64ToBytes 转bytes
func Uint64ToBytes(n uint64) []byte {
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, app.NetEndian, &n)
	return bytesBuffer.Bytes()
}

//Int64ToBytes 转bytes
func Int64ToBytes(n int64) []byte {
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, app.NetEndian, &n)
	return bytesBuffer.Bytes()
}

//CompressBytes 压缩字节
func CompressBytes(data []byte) ([]byte, error) {

	var buf bytes.Buffer

	writer := zlib.NewWriter(&buf)

	_, err := writer.Write(data)
	if err != nil {
		return nil, err
	}
	writer.Close()

	return buf.Bytes(), nil
}

//DecompressBytes 解压字节
func DecompressBytes(data []byte) ([]byte, error) {

	reader, err := zlib.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	defer reader.Close()

	return ioutil.ReadAll(reader)
}
