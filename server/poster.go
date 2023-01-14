package main

import (
	"errors"

	"github.com/huoyijie/GoChat/lib"
	"google.golang.org/protobuf/proto"
)

// 转换同步响应类型
func syncResponseToKind(m proto.Message) (kind lib.PackKind, err error) {
	switch m.(type) {
	case *lib.TokenRes, *lib.UsersRes:
		kind = lib.PackKind_RES
	default:
		err = errors.New("invalid kind of packet")
	}
	return
}

// 实现 post 接口
type poster struct {
	packChan chan<- *lib.Packet
}

func newPoster(packChan chan<- *lib.Packet) *poster {
	return &poster{packChan}
}

// Handle implements lib.Post
func (p *poster) Handle(req, res proto.Message) (err error) {
	pack, ok := req.(*lib.Packet)
	if !ok {
		return errors.New("invalid request")
	}

	kind, err := syncResponseToKind(res)
	if err != nil {
		return
	}

	bytes, err := lib.Marshal(res)
	if err != nil {
		return
	}

	p.packChan <- &lib.Packet{
		Id:   pack.Id,
		Kind: kind,
		Data: bytes,
	}

	return
}

// Send implements lib.Post
func (p *poster) Send(res proto.Message) (err error) {
	var kind lib.PackKind
	switch res.(type) {
	case *lib.Pong:
		kind = lib.PackKind_PONG
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

	p.packChan <- &lib.Packet{
		Kind: kind,
		Data: bytes,
	}

	return
}

// Close implements lib.Post
func (p *poster) Close() {
	close(p.packChan)
}

var _ lib.Post = (*poster)(nil)
