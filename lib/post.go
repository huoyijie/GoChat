package lib

import "google.golang.org/protobuf/proto"

// 定义处理 packet 接口
type Post interface {
	// 处理同步请求
	Handle(req, res proto.Message) error

	// 发送非同步请求
	Send(req proto.Message) error
}
