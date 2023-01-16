package main

import (
	"github.com/huoyijie/GoChat/lib"
	"google.golang.org/protobuf/proto"
)

// 处理 biz_ping_t 请求
type biz_ping_t struct {
	biz_base_t
}

func initialPing(base biz_base_t) *biz_ping_t {
	return &biz_ping_t{base}
}

func (p *biz_ping_t) do(req proto.Message, accId *uint64, accUN *string) error {
	pack, err := p.toPacket(req)
	if err != nil {
		return err
	}

	ping := &lib.Ping{}
	if err := lib.Unmarshal(pack.Data, ping); err != nil {
		return err
	}
	// log.Println(string(ping.Payload))

	return p.poster.Send(&lib.Pong{Payload: []byte("宝塔镇河妖")})
}

var _ biz_i = (*biz_ping_t)(nil)
