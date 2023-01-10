package main

import (
	"errors"

	"github.com/huoyijie/GoChat/lib"
	"google.golang.org/protobuf/proto"
)

func syncResponseToKind(m proto.Message) (kind lib.PackKind, err error) {
	switch m.(type) {
	case *lib.TokenRes, *lib.UsersRes:
		kind = lib.PackKind_RES
	default:
		err = errors.New("invalid kind of packet")
	}
	return
}

func handlePacket(packChan chan<- *lib.Packet, pack *lib.Packet, res proto.Message) (err error) {
	kind, err := syncResponseToKind(res)
	if err != nil {
		return
	}

	bytes, err := lib.Marshal(res)
	if err != nil {
		return
	}

	packChan <- &lib.Packet{
		Id:   pack.Id,
		Kind: kind,
		Data: bytes,
	}
	return
}

func sendPacket(packChan chan<- *lib.Packet, res proto.Message) (err error) {
	var kind lib.PackKind
	switch res.(type) {
	case *lib.Msg:
		kind = lib.PackKind_MSG
	case *lib.ErrRes:
		kind = lib.PackKind_ERR
	default:
		return errors.New("invalid kind of packet")
	}

	bytes, err := lib.Marshal(res)
	if err != nil {
		return
	}

	packChan <- &lib.Packet{
		Kind: kind,
		Data: bytes,
	}
	return
}
