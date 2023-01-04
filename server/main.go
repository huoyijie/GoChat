package main

import (
	"fmt"
	"net"
	"sync"

	"github.com/bwmarrin/snowflake"
	"github.com/huoyijie/GoChat/lib"
)

// 存储当前所有客户端连接
var sockets = make(map[snowflake.ID]*lib.Socket)

// 多个协程并发读写 sockets 时，需要使用读写锁
var lock sync.RWMutex

// 写锁
func wSockets(wSockets func()) {
	lock.Lock()
	defer lock.Unlock()
	wSockets()
}

// 读锁
func rSockets(rSockets func()) {
	lock.RLock()
	defer lock.RUnlock()
	rSockets()
}

func main() {
	// 启动单独协程，监听 ctrl+c 或 kill 信号，收到信号结束进程
	go lib.SignalHandler()

	// tcp 监听地址 0.0.0.0:8888
	addr := ":8888"

	// tcp 监听
	ln, err := net.Listen("tcp", addr)

	// tcp 监听遇到错误退出进程
	lib.FatalNotNil(err)
	// 输出日志
	lib.LogMessage("Listening on", addr)

	// 创建 snowflake Node
	node, err := snowflake.NewNode(1)
	lib.FatalNotNil(err)

	// 循环接受客户端连接
	for {
		// 每当有客户端连接时，ln.Accept 会返回新的连接 conn
		conn, err := ln.Accept()
		// 如果接受的新连接遇到错误，则退出进程
		lib.FatalNotNil(err)

		// 生成新 ID
		id := node.Generate()
		packChan := make(chan *lib.Packet)
		// 保存新连接
		wSockets(func() {
			sockets[id] = &lib.Socket{Id: id, PackChan: packChan}
		})

		go lib.SendTo(conn, packChan)

		// 为每个客户端启动一个协程，读取客户端发送的消息并转发
		go lib.RecvFrom(
			conn,
			// todo return error and close conn
			func(pack *lib.Packet) {
				switch pack.Kind {
				case lib.PackKind_MSG:
					rSockets(func() {
						for k, v := range sockets {
							// 向其他所有客户端(除了自己)转发消息
							if k != id {
								v.PackChan <- &lib.Packet{
									Kind: pack.Kind,
									Data: pack.Data,
								}
							}
						}
					})
				case lib.PackKind_SIGNUP:
					signup := &lib.Signup{}
					if err := lib.Unmarshal(pack.Data, signup); err == nil {
						fmt.Println(signup.Auth)
					}
				}
			},
			func() {
				// 从当前方法返回时，关闭连接
				conn.Close()
				// 删除连接
				wSockets(func() {
					delete(sockets, id)
				})
			})
	}
}
