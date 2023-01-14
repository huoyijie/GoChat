package main

import (
	"github.com/huoyijie/GoChat/lib"
	"google.golang.org/protobuf/proto"
)

// 处理 ping 请求
type ping struct {
	base
}

func initialPing(base base) *ping {
	return &ping{base}
}

func (p *ping) do(req proto.Message, accId *uint64, accUN *string) error {
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

var _ biz = (*ping)(nil)
