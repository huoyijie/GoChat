package main

import (
	"errors"

	"github.com/huoyijie/GoChat/lib"
	"google.golang.org/protobuf/proto"
)

// 定义服务器后台业务逻辑接口
type biz_i interface {
	// req: 当前请求信息
	//
	// accId: 当前登录用户 id
	//
	// accUN: 当前登录用户 username
	do(req proto.Message, accId *uint64, accUN *string) error
}

// 后台业务逻辑对象可嵌入 biz_base_t
type biz_base_t struct {
	poster  lib.Post
	storage *storage_t
}

func initialBase(poster lib.Post, storage *storage_t) biz_base_t {
	return biz_base_t{
		poster,
		storage,
	}
}

// 生成 token 并向客户端发送 TokenRes packet
func (b *biz_base_t) handleAuth(pack *lib.Packet, account *Account, accId *uint64, accUN *string) error {
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
func (b *biz_base_t) unmarshal(pack *lib.Packet, req proto.Message) error {
	if err := lib.Unmarshal(pack.Data, req); err != nil {
		return b.poster.Handle(pack, &lib.TokenRes{Code: lib.Err_Unmarshal.Val()})
	}
	return nil
}

// 把 req 对象转换为 packet
func (b *biz_base_t) toPacket(req proto.Message) (pack *lib.Packet, err error) {
	pack, ok := req.(*lib.Packet)
	if !ok {
		err = errors.New("invalid request")
	}
	return
}

func (b *biz_base_t) close() {
	b.poster.Close()
}
