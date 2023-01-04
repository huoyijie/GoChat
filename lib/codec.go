package lib

import (
	b "bytes"
	"encoding/binary"
	"errors"

	"google.golang.org/protobuf/proto"
)

const PackLenField = 4

func Uint32ToBytes(n uint32) (bytes []byte) {
	bytes = make([]byte, PackLenField)
	binary.BigEndian.PutUint32(bytes, n)
	return
}

func BytesToUint32(bytes []byte) uint32 {
	return binary.BigEndian.Uint32(bytes)
}

func SplitFunc(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if len(data) >= PackLenField {
		length := BytesToUint32(data[:PackLenField])

		if uint32(len(data)-PackLenField) >= length {
			advance = int(PackLenField + length)
			token = make([]byte, length)
			copy(token, data[PackLenField:advance])
		}
	}

	if atEOF && len(data[advance:]) > 0 {
		err = errors.New("遇到 EOF，但数据不完整，未能解析出数据包！")
	}

	return
}

func MarshalPack(pack proto.Message) (bytes []byte, err error) {
	packBytes, err := proto.Marshal(pack)
	if err != nil {
		return
	}

	packLenFieldBytes := Uint32ToBytes(uint32(len(packBytes)))
	bytes = b.Join(
		[][]byte{
			packLenFieldBytes,
			packBytes,
		},
		[]byte{},
	)
	return
}

func Marshal(m proto.Message) (bytes []byte, err error) {
	bytes, err = proto.Marshal(m)
	return
}

func Unmarshal(bytes []byte, m proto.Message) (err error) {
	err = proto.Unmarshal(bytes, m)
	return
}
