package main

import (
	"github.com/huoyijie/GoChat/lib"
	"google.golang.org/protobuf/proto"
)

// 转发消息逻辑
type fw_msg struct {
	base
}

func initialFwMsg(base base) *fw_msg {
	return &fw_msg{base}
}

func (fwm *fw_msg) do(req proto.Message, accId *uint64, accUN *string) error {
	if len(*accUN) > 0 { // accUN 用户已登录客户端
		// 查询所有发送给 accUN 的未读消息
		msgList, _ := fwm.storage.GetMsgList(*accUN)
		for i := range msgList {
			// 转发消息
			fwm.poster.Send(&lib.Msg{
				Id:   msgList[i].Id,
				Kind: lib.MsgKind(msgList[i].Kind),
				From: msgList[i].From,
				To:   msgList[i].To,
				Data: msgList[i].Data,
			})
		}
	}
	return nil
}

var _ biz = (*fw_msg)(nil)
