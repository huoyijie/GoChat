package main

import (
	"time"

	"github.com/huoyijie/GoChat/lib"
)

// 封装服务器请求，会向服务器发送 packet
type request_t struct {
	pack     *lib.Packet
	c        chan *response_t
	deadline time.Time
}

// 创建服务器请求对象
func newRequest(pack *lib.Packet) (request *request_t) {
	request = &request_t{pack: pack}
	// 同步请求发送后，可通过 request.c channel 获取响应
	if sync := pack.Kind > lib.PackKind_PING; sync {
		request.c = make(chan *response_t, 1)
	}
	return
}

// 判断当前请求是否为同步请求
func (request *request_t) sync() bool {
	return request.c != nil
}

// 封装服务器响应
type response_t struct {
	pack *lib.Packet
}

// 判断服务器是否正常返回了数据
func (response *response_t) ok() bool {
	return response.pack != nil
}
