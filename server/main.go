package main

import (
	"bufio"
	"errors"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/huoyijie/GoChat/lib"
)

// 信号监听处理器
func signalHandler() {
	// 创建信号 channel
	sigchan := make(chan os.Signal, 1)

	// 注册要监听哪些信号
	signal.Notify(sigchan, os.Interrupt)    // ctrl+c
	signal.Notify(sigchan, syscall.SIGTERM) // kill

	// 一直阻塞，直到收到信号，恢复执行并退出进程
	<-sigchan
	// 退出进程
	defer os.Exit(0)
}

// 把来自 packChan 的 packet 都发送到 conn
func sendTo(conn net.Conn, b base, interval *time.Ticker, packChan <-chan *lib.Packet, accId *uint64, accUN *string) {
	// 从当前方法返回后，断开连接，清理资源等
	defer conn.Close()
	defer interval.Stop()

	var (
		fw_msg biz = initialFwMsg(b)
		pid    uint64
	)

	for {
		select {

		// 发送 packet 到服务器
		case pack, ok := <-packChan:
			if !ok { // recvFrom 协程已退出，需要退出当前协程
				return
			}

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

		// 读取并转发所有发送给当前用户的未读消息
		case <-interval.C:
			fw_msg.do(nil, accId, accUN)
		}
	}
}

// 根据 kind 返回对应的后台处理逻辑 biz
func kindToBiz(kind lib.PackKind, b base, node *snowflake.Node) (biz biz, err error) {
	switch kind {
	case lib.PackKind_SIGNUP:
		biz = initialSignup(b)
	case lib.PackKind_SIGNIN:
		biz = initialSignin(b)
	case lib.PackKind_TOKEN:
		biz = initialValToken(b)
	case lib.PackKind_USERS:
		biz = initialUsers(b)
	case lib.PackKind_MSG:
		biz = initialRecvMsg(b, node)
	default:
		err = errors.New("invalid kind of packet")
	}
	return
}

// 读取并处理客户端发送的 packet
func recvFrom(conn net.Conn, b base, interval *time.Ticker, accId *uint64, accUN *string, node *snowflake.Node) {
	defer b.close()
	defer interval.Stop()

	// 设置如何处理接收到的字节流，SplitFunc 会根据 packet 开头 length 把字节流分割为消息流
	scanner := bufio.NewScanner(conn)
	scanner.Split(lib.SplitFunc)

	// 循环解析消息，每当解析出一条消息后，scan() 返回 true
	for scanner.Scan() {
		// 把 scanner 解析出的消息字节 slice 解析为 Pack
		pack := &lib.Packet{}
		if err := lib.Unmarshal(scanner.Bytes(), pack); err != nil {
			return
		}

		// 获取 packet 处理逻辑
		biz, err := kindToBiz(pack.Kind, b, node)
		if err != nil {
			return
		}

		// 执行 packet 处理逻辑
		if err := biz.do(pack, accId, accUN); err != nil {
			return
		}
	}
}

func handleConn(conn net.Conn, storage *Storage, node *snowflake.Node) {
	// 当前连接所登录的用户
	var (
		accId uint64
		accUN string
	)

	// 通过该 channel 可向当前连接发送 packet
	packChan := make(chan *lib.Packet)
	var poster lib.Post = newPoster(packChan)
	base := initialBase(poster, storage)
	// 间隔 100ms 检查是否有新消息
	interval := time.NewTicker(100 * time.Millisecond)

	// 为每个客户端启动一个协程，读取并处理客户端发送的 packet
	go recvFrom(conn, base, interval, &accId, &accUN, node)

	// 为每个客户端启动一个协程，把来自 packChan 的 packet 都发送到 conn
	sendTo(conn, base, interval, packChan, &accId, &accUN)
}

func main() {
	// 启动单独协程，监听 ctrl+c 或 kill 信号，收到信号结束进程
	go signalHandler()

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

		// 启动新协程处理当前连接
		go handleConn(conn, storage, node)
	}
}
