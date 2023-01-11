package main

import (
	"errors"

	"github.com/huoyijie/GoChat/lib"
	"google.golang.org/protobuf/proto"
)

// 定义服务器后台业务逻辑接口
type biz interface {
	// req: 当前请求信息
	//
	// accId: 当前登录用户 id
	//
	// accUN: 当前登录用户 username
	do(req proto.Message, accId *uint64, accUN *string) error
}

// 后台业务逻辑对象可嵌入 base
type base struct {
	poster  lib.Post
	storage *Storage
}

func initialBase(poster lib.Post, storage *Storage) base {
	return base{
		poster,
		storage,
	}
}

// 生成 token 并向客户端发送 TokenRes packet
func (b *base) handleAuth(pack *lib.Packet, account *Account, accId *uint64, accUN *string) error {
	token, err := GenerateToken(account.Id)
	if err != nil {
		return b.poster.Handle(pack, &lib.TokenRes{Code: lib.Err_Gen_Token.Val()})
	}

	err = b.poster.Handle(pack, &lib.TokenRes{Id: account.Id, Username: account.Username, Token: token})
	if err != nil {
		return err
	}
	*accId = account.Id
	*accUN = account.Username
	return nil
}

// 反序列化请求对象
func (b *base) unmarshal(pack *lib.Packet, req proto.Message) error {
	if err := lib.Unmarshal(pack.Data, req); err != nil {
		return b.poster.Handle(pack, &lib.TokenRes{Code: lib.Err_Unmarshal.Val()})
	}
	return nil
}

// 把 req 对象转换为 packet
func (b *base) toPacket(req proto.Message) (pack *lib.Packet, err error) {
	pack, ok := req.(*lib.Packet)
	if !ok {
		err = errors.New("invalid request")
	}
	return
}
