package main

import (
	"time"

	"github.com/huoyijie/GoChat/lib"
)

type request_t struct {
	pack     *lib.Packet
	c        chan *response_t
	deadline time.Time
}

func (request *request_t) init(pack *lib.Packet) *request_t {
	request.pack = pack
	switch {
	case pack.Kind > lib.PackKind_PING:
		request.c = make(chan *response_t, 1)
	}
	return request
}

func (request *request_t) sync() bool {
	return request.c != nil
}

type response_t struct {
	pack *lib.Packet
}

func (response *response_t) ok() bool {
	return response.pack != nil
}