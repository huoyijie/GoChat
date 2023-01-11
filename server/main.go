package main

import (
	"errors"
	"net"
	"path/filepath"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/huoyijie/GoChat/lib"
)

// 把来自 packChan 的 packet 都发送到 conn
func sendTo(conn net.Conn, packChan <-chan *lib.Packet) {
	var pid uint64
	for pack := range packChan {
		if pack.Id == 0 {
			pid++
			pack.Id = pid
		}

		bytes, err := lib.MarshalPack(pack)
		if err != nil {
			return
		}

		_, err = conn.Write(bytes)
		if err != nil {
			return
		}
	}
}

// 读取并处理客户端发送的 packet
func recvFrom(conn net.Conn, b base, accId *uint64, accUN *string, node *snowflake.Node) {
	var (
		signup    biz = initialSignup(b)
		signin    biz = initialSignin(b)
		val_token biz = initialValToken(b)
		users     biz = initialUsers(b)
		recv_msg  biz = initialRecvMsg(b, node)
	)
	lib.RecvFrom(
		conn,
		func(pack *lib.Packet) (err error) {
			var biz biz
			switch pack.Kind {
			case lib.PackKind_SIGNUP:
				biz = signup
			case lib.PackKind_SIGNIN:
				biz = signin
			case lib.PackKind_TOKEN:
				biz = val_token
			case lib.PackKind_USERS:
				biz = users
			case lib.PackKind_MSG:
				biz = recv_msg
			default:
				return errors.New("invalid kind of packet")
			}
			err = biz.do(pack, accId, accUN)
			return
		})
}

// 读取并转发所有发送给 accUN 用户的未读消息
func forwardMsgs(b base, accId *uint64, accUN *string) {
	// 间隔 100ms 检查是否有新消息
	interval := time.NewTicker(100 * time.Millisecond)
	defer interval.Stop()

	var fw_msg biz = initialFwMsg(b)

	for range interval.C {
		fw_msg.do(nil, accId, accUN)
	}
}

func main() {
	// 启动单独协程，监听 ctrl+c 或 kill 信号，收到信号结束进程
	go lib.SignalHandler()

	// 初始化存储
	storage, err := new(Storage).Init(filepath.Join(lib.WorkDir, "server.db"))
	lib.FatalNotNil(err)

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

		// 当前连接所登录的用户
		var (
			accId uint64
			accUN string
		)

		// 通过该 channel 可向当前连接发送 packet
		packChan := make(chan *lib.Packet)
		poster := newPoster(packChan)
		base := initialBase(poster, storage)

		// 为每个客户端启动一个协程，把来自 packChan 的 packet 都发送到 conn
		go sendTo(conn, packChan)

		// 为每个客户端启动一个协程，读取并处理客户端发送的 packet
		go recvFrom(conn, base, &accId, &accUN, node)

		// 为每个客户端启动一个协程，读取并转发所有发送给当前用户的未读消息
		go forwardMsgs(base, &accId, &accUN)
	}
}
