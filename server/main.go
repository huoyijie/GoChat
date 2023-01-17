package main

import (
	"bufio"
	"errors"
	"log"
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
	sigChan := make(chan os.Signal, 1)

	// 注册要监听哪些信号
	signal.Notify(sigChan, os.Interrupt)    // ctrl+c
	signal.Notify(sigChan, syscall.SIGTERM) // kill

	// 一直阻塞，直到收到信号，恢复执行并退出进程
	<-sigChan
	// 退出进程
	defer os.Exit(0)
}

// 把来自 packChan 的 packet 都发送到 conn
func sendTo(conn net.Conn, packChan <-chan *lib.Packet, c <-chan *lib.Push, accId *uint64, accUN *string, storage *storage_t) {
	// 间隔 100ms 检查是否有新消息
	interval := time.NewTicker(100 * time.Millisecond)
	defer interval.Stop()
	var pid uint64

	var sendPack = func(pack *lib.Packet) (err error) {
		if pack.Id == 0 {
			pid++
			pack.Id = pid
		}

		bytes, err := lib.MarshalPack(pack)
		if err != nil {
			return err
		}

		_, err = conn.Write(bytes)
		return
	}

loop:
	for {
		select {

		// 发送 packet 到客户端
		case pack, ok := <-packChan:
			if !ok { // recvFrom 协程已退出，需要退出当前协程
				return
			}

			if err := sendPack(pack); err != nil {
				log.Println(err)
				return
			}

		// 发送 push 到客户端
		case push := <-c:
			// 上下线提醒不用发给自己
			if push.Kind == lib.PushKind_ONLINE {
				online := &lib.Online{}
				if err := lib.Unmarshal(push.Data, online); err != nil {
					log.Println(err)
					return
				}
				if online.Username == *accUN {
					continue loop
				}
			}

			bytes, err := lib.Marshal(push)
			if err != nil {
				log.Println(err)
				return
			}

			pack := &lib.Packet{
				Kind: lib.PackKind_PUSH,
				Data: bytes,
			}
			if err := sendPack(pack); err != nil {
				log.Println(err)
				return
			}

		// 读取并转发所有发送给当前用户的未读消息
		case <-interval.C:
			if len(*accUN) > 0 { // accUN 用户已登录客户端
				// 查询所有发送给 accUN 的未读消息
				msgList, _ := storage.GetMsgList(*accUN)
				for i := range msgList {
					msg := &lib.Msg{
						Id:   msgList[i].Id,
						Kind: lib.MsgKind(msgList[i].Kind),
						From: msgList[i].From,
						To:   msgList[i].To,
						Data: msgList[i].Data,
					}

					bytes, err := lib.Marshal(msg)
					if err != nil {
						log.Println(err)
						return
					}

					// 转发消息
					if err := sendPack(&lib.Packet{
						Kind: lib.PackKind_MSG,
						Data: bytes,
					}); err != nil {
						log.Println(err)
						return
					}
				}
			}
		}
	}
}

// 根据 kind 返回对应的后台处理逻辑 biz
func kindToBiz(kind lib.PackKind, b biz_base_t, node *snowflake.Node) (biz biz_i, err error) {
	switch kind {
	case lib.PackKind_PING:
		biz = initialPing(b)
	case lib.PackKind_SIGNUP:
		biz = initialSignup(b)
	case lib.PackKind_SIGNIN:
		biz = initialSignin(b)
	case lib.PackKind_TOKEN:
		biz = initialValToken(b)
	case lib.PackKind_SIGNOUT:
		biz = initialSignout(b)
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
func recvFrom(conn net.Conn, b biz_base_t, accId *uint64, accUN *string, node *snowflake.Node) {
	defer b.close()

	// 设置如何处理接收到的字节流，SplitFunc 会根据 packet 开头 length 把字节流分割为消息流
	scanner := bufio.NewScanner(conn)
	scanner.Split(lib.SplitFunc)

	// 循环解析消息，每当解析出一条消息后，scan() 返回 true
	for scanner.Scan() {
		// 把 scanner 解析出的消息字节 slice 解析为 Pack
		pack := &lib.Packet{}
		if err := lib.Unmarshal(scanner.Bytes(), pack); err != nil {
			log.Println(err)
			return
		}

		// 获取 packet 处理逻辑
		biz, err := kindToBiz(pack.Kind, b, node)
		if err != nil {
			log.Println(err)
			return
		}

		// 执行 packet 处理逻辑
		if err := biz.do(pack, accId, accUN); err != nil {
			log.Println(err)
			return
		}
	}
}

func handleConn(conn net.Conn, sid uint64, eventChan chan<- event_i, pushChan chan<- *lib.Push, storage *storage_t, node *snowflake.Node) {
	// 从当前方法返回后，断开连接，清理资源等
	defer conn.Close()

	// 当前连接所登录的用户
	var (
		accId uint64
		accUN string
	)

	// 断开连接后，更新用户在线状态
	defer func() {
		if accId > 0 {
			storage.UpdateOnline(accId, false)

			// 下线事件
			eventChan <- &e_offline_t{sid}

			// 下线提醒
			bytes, err := lib.Marshal(&lib.Online{Kind: lib.OnlineKind_OFF, Username: accUN})
			lib.FatalNotNil(err)

			pushChan <- &lib.Push{
				Kind: lib.PushKind_ONLINE,
				Data: bytes,
			}
		}
	}()

	// 通过该 channel 可向当前连接发送 packet
	packChan := make(chan *lib.Packet, 1024)
	var poster lib.Post = newPoster(packChan)
	base := initialBase(sid, poster, eventChan, pushChan, storage)

	// 为每个客户端启动一个协程，读取并处理客户端发送的 packet
	go recvFrom(conn, base, &accId, &accUN, node)

	// 为每个客户端启动一个协程，把来自 packChan 的 packet 都发送到 conn
	sendTo(conn, packChan, base.c, &accId, &accUN, storage)
}

func main() {
	// 启动单独协程，监听 ctrl+c 或 kill 信号，收到信号结束进程
	go signalHandler()

	// 初始化存储
	storage, err := new(storage_t).Init(filepath.Join(lib.WorkDir, "server.db"))
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

	eventChan := make(chan event_i, 1024)
	pushChan := make(chan *lib.Push, 1024)
	// 开启独立协程处理 push
	go handlePush(eventChan, pushChan)

	var sid uint64
	// 循环接受客户端连接
	for {
		// 每当有客户端连接时，ln.Accept 会返回新的连接 conn
		conn, err := ln.Accept()
		// 如果接受的新连接遇到错误，则退出进程
		lib.FatalNotNil(err)

		sid++

		// 启动新协程处理当前连接
		go handleConn(conn, sid, eventChan, pushChan, storage, node)
	}
}
