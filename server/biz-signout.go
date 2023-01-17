package main

import (
	"github.com/huoyijie/GoChat/lib"
	"google.golang.org/protobuf/proto"
)

// 处理登出请求
type biz_signout_t struct {
	biz_base_t
}

func initialSignout(base biz_base_t) *biz_signout_t {
	return &biz_signout_t{base}
}

func (s *biz_signout_t) do(req proto.Message, accId *uint64, accUN *string) error {
	s.storage.UpdateOnline(*accId, false)

	// 下线事件
	s.eventChan <- &offline_t{s.sid}

	// 下线提醒
	bytes, err := lib.Marshal(&lib.Online{Kind: lib.OnlineKind_OFF, Username: *accUN})
	if err != nil {
		return err
	}
	s.pushChan <- &lib.Push{
		Kind: lib.PushKind_ONLINE,
		Data: bytes,
	}

	*accId = 0
	*accUN = ""

	return s.poster.Handle(req, &lib.SignoutRes{})
}

var _ biz_i = (*biz_signout_t)(nil)
