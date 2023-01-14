package main

import (
	"errors"
	"fmt"

	"github.com/huoyijie/GoChat/lib"
	"google.golang.org/protobuf/proto"
)

// 转换同步请求类型
func syncRequestToKind(m proto.Message) (kind lib.PackKind, err error) {
	switch m.(type) {
	case *lib.Signup:
		kind = lib.PackKind_SIGNUP
	case *lib.Signin:
		kind = lib.PackKind_SIGNIN
	case *lib.Token:
		kind = lib.PackKind_TOKEN
	case *lib.Users:
		kind = lib.PackKind_USERS
	default:
		err = errors.New("invalid kind of packet")
	}
	return
}

// 实现 post 接口
type poster struct {
	reqChan chan<- *request_t
}

func newPoster(reqChan chan<- *request_t) *poster {
	return &poster{reqChan}
}

// Handle implements lib.Post
func (p *poster) Handle(req, res proto.Message) (err error) {
	kind, err := syncRequestToKind(req)
	if err != nil { // 转换同步请求类型
		return
	}

	bytes, err := lib.Marshal(req)
	if err != nil { // 序列化请求错误
		return
	}

	request := newRequest(&lib.Packet{Kind: kind, Data: bytes})

	p.reqChan <- request
	response := <-request.c
	if !response.ok() { // 同步请求超时
		err = fmt.Errorf("%s 请求超时", kind)
		return
	}

	err = lib.Unmarshal(response.pack.Data, res)
	return
}

// Send implements lib.Post
func (p *poster) Send(req proto.Message) (err error) {
	var kind lib.PackKind
	switch req.(type) {
	case *lib.Msg:
		kind = lib.PackKind_MSG
	default:
		return errors.New("invalid kind of packet")
	}

	bytes, err := lib.Marshal(req)
	if err != nil {
		return
	}

	request := newRequest(&lib.Packet{Kind: kind, Data: bytes})

	p.reqChan <- request
	return
}

// Close implements lib.Post
func (p *poster) Close() {
	close(p.reqChan)
}

var _ lib.Post = (*poster)(nil)
