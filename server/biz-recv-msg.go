package main

import (
	"github.com/bwmarrin/snowflake"
	"github.com/huoyijie/GoChat/lib"
	"google.golang.org/protobuf/proto"
)

// 处理发送消息请求
type recv_msg struct {
	base
	node *snowflake.Node
}

func initialRecvMsg(base base, node *snowflake.Node) *recv_msg {
	return &recv_msg{base, node}
}

func (rm *recv_msg) do(req proto.Message, accId *uint64, accUN *string) error {
	pack, err := rm.toPacket(req)
	if err != nil {
		return err
	}

	if len(*accUN) == 0 {
		return rm.poster.Send(&lib.ErrRes{Code: lib.Err_Forbidden.Val()})
	}

	msg := &lib.Msg{}
	if err := lib.Unmarshal(pack.Data, msg); err != nil {
		return rm.poster.Send(&lib.ErrRes{Code: lib.Err_Unmarshal.Val()})
	}

	return rm.storage.NewMsg(&Message{
		// 生成消息 ID
		Id:   int64(rm.node.Generate()),
		Kind: uint32(lib.MsgKind_TEXT),
		From: msg.From,
		To:   msg.To,
		Data: msg.Data,
	})
}

var _ biz = (*recv_msg)(nil)
