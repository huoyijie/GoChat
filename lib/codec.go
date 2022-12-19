package lib

import (
	b "bytes"
	"encoding/binary"
	"errors"
	"net"

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

func SendMsg(conn net.Conn, msg proto.Message) (err error) {
	msgBytes, err := proto.Marshal(msg)
	if err != nil {
		return
	}

	err = sendPack(conn, &Packet{
		Id:   1,
		Kind: 1,
		Data: msgBytes,
	})

	return
}

func RecvMsg(bytes []byte) (msg proto.Message, err error) {
	pack, err := recvPack(bytes)
	if err != nil {
		return
	}

	msg = &Msg{}
	err = proto.Unmarshal(pack.Data, msg)
	return
}

func sendPack(conn net.Conn, pack proto.Message) (err error) {
	packBytes, err := proto.Marshal(pack)
	if err != nil {
		return
	}

	packLenFieldBytes := Uint32ToBytes(uint32(len(packBytes)))
	bytes := b.Join(
		[][]byte{
			packLenFieldBytes,
			packBytes,
		},
		[]byte{},
	)

	_, err = conn.Write(bytes)
	return
}

func recvPack(bytes []byte) (pack *Packet, err error) {
	pack = &Packet{}
	err = proto.Unmarshal(bytes, pack)
	return
}