package main

import (
	"github.com/huoyijie/GoChat/lib"
)

type event_i interface{}

type online_t struct {
	sid uint64
	c   chan<- *lib.Push
}

type offline_t struct {
	sid uint64
}

func handlePush(eventChan <-chan event_i, pushChan <-chan *lib.Push) {
	sessions := make(map[uint64]chan<- *lib.Push)
	for {
		select {
		case e := <-eventChan:
			switch e := e.(type) {
			case *online_t:
				sessions[e.sid] = e.c
			case *offline_t:
				delete(sessions, e.sid)
			}
		case push := <-pushChan:
			for _, c := range sessions {
				c <- push
			}
		}
	}
}
