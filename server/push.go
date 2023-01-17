package main

import (
	"github.com/huoyijie/GoChat/lib"
)

// 定义事件接口
type event_i interface{}

// 上线事件
type e_online_t struct {
	sid uint64
	c   chan<- *lib.Push
}

// 下线事件
type e_offline_t struct {
	sid uint64
}

// 维护客户端 sessions，接收并处理客户端上下线事件，接收并转发 push 到客户端
func handlePush(eventChan <-chan event_i, pushChan <-chan *lib.Push) {
	sessions := make(map[uint64]chan<- *lib.Push)
	for {
		select {
		case e := <-eventChan:
			switch e := e.(type) {
			case *e_online_t:
				sessions[e.sid] = e.c
			case *e_offline_t:
				delete(sessions, e.sid)
			}
		case push := <-pushChan:
			for _, c := range sessions {
				c <- push
			}
		}
	}
}
