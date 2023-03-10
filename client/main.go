package main

import (
	"bufio"
	"encoding/base64"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/huoyijie/GoChat/lib"
)

// 向服务器发送 packet。如果是同步请求，会通过 request.c 返回服务器响应数据，同时也会检查同步请求是否已超时。
func sendTo(conn net.Conn, reqChan <-chan *request_t, resChan <-chan *response_t) (quit bool) {
	// 从当前方法返回后，断开连接，清理资源等
	defer conn.Close()

	// 分配 packet.id
	var id uint64
	// 登记所有的同步请求，并等待响应
	requests := make(map[uint64]*request_t)
	// 同步请求超时检查间隔 50ms
	timeoutTicker := time.NewTicker(50 * time.Millisecond)
	defer timeoutTicker.Stop()
	pingTicker := time.NewTicker(20 * time.Second)
	defer pingTicker.Stop()

	var sendPack = func(pack *lib.Packet) (err error) {
		id++
		pack.Id = id

		bytes, err := lib.MarshalPack(pack)
		if err != nil { // 序列化 packet 错误
			return
		}

		_, err = conn.Write(bytes)
		return
	}

	for {
		select {

		// 有服务器请求进来
		case request, ok := <-reqChan:
			if !ok { // renderUI 协程已退出，需要退出进程
				quit = true
				return
			}

			if err := sendPack(request.pack); err != nil { // 发送字节数据错误
				return
			}

			if request.sync() { // 登记同步请求
				// 如果 5s 内服务器没有返回，则超时
				request.deadline = time.Now().Add(5 * time.Second)
				requests[request.pack.Id] = request
			}

		// 通过 resChan 接收从 recvFrom 协程发送过来的响应对象
		case response, ok := <-resChan:
			if !ok { // recvFrom 协程已退出，连接已断开，需要重新连接和启动新协程
				return
			}

			if request, found := requests[response.pack.Id]; found {
				// 通过 request.c 返回服务器响应数据
				request.c <- response
				// 删除登记的同步请求
				delete(requests, response.pack.Id)
			}

		// 每隔 50ms 检查是否有超时的同步请求
		case <-timeoutTicker.C:
			for id, request := range requests {
				if time.Now().After(request.deadline) {
					// 请求超时，返回空响应对象
					request.c <- newResponse(nil)
					// 删除登记的同步请求
					delete(requests, id)
				}
			}

		// ping
		case <-pingTicker.C:
			bytes, err := lib.Marshal(&lib.Ping{Payload: []byte("天王盖地虎")})
			if err != nil {
				return
			}

			if err := sendPack(&lib.Packet{
				Kind: lib.PackKind_PING,
				Data: bytes,
			}); err != nil {
				return
			}
		}
	}
}

// 从服务器接收 packet 的处理函数
func handlePack(pack *lib.Packet, resChan chan<- *response_t, storage *storage_t) (err error) {
	switch pack.Kind {

	// pong
	case lib.PackKind_PONG:
		pong := &lib.Pong{}
		if err := lib.Unmarshal(pack.Data, pong); err != nil {
			return err
		}
		// log.Println(string(pong.Payload))

	// 当前连接用户收到服务器 push
	case lib.PackKind_PUSH:
		push := &lib.Push{}
		err = lib.Unmarshal(pack.Data, push)
		if err != nil {
			return
		}
		// 新 push 写入本地存储
		storage.NewPush(&Push{Kind: int32(push.Kind), Data: push.Data})

	// 当前连接的登录用户收到新未读消息
	case lib.PackKind_MSG:
		msg := &lib.Msg{}
		err = lib.Unmarshal(pack.Data, msg)
		if err != nil {
			return
		}
		// 新消息写入本地存储
		storage.NewMsg(&Message{
			Id:   msg.Id,
			Kind: int32(msg.Kind),
			From: msg.From,
			Data: msg.Data,
		})

	// 当前连接遇到系统异常，退出进程
	case lib.PackKind_ERR:
		errRes := &lib.ErrRes{}
		err = lib.Unmarshal(pack.Data, errRes)
		if err != nil {
			return
		}
		err = fmt.Errorf("系统异常: %d", errRes.Code)

	// 收到同步请求的响应
	case lib.PackKind_RES:
		// 当前是在 recvFrom 协程里，需要把 packet 封装为 response_t 对象，并通过 resChan channel 发到 sendTo 协程
		resChan <- newResponse(pack)
	}
	return
}

// 从服务器接收 packet 并进行处理
func recvFrom(conn net.Conn, resChan chan<- *response_t, storage *storage_t) {
	// 协程退出前关闭 channel
	defer close(resChan)

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

		// 执行 packet 处理逻辑
		if err := handlePack(pack, resChan, storage); err != nil {
			return
		}
	}
}

// 验证 token 是否有效
func validateToken(poster lib.Post, storage *storage_t) (tokenRes *lib.TokenRes, err error) {
	kv, err := storage.GetValue("token")
	if err != nil { // 未登录过
		return
	}

	token, err := base64.StdEncoding.DecodeString(kv.Value)
	if err != nil { // token 解析错误
		return
	}

	tokenRes = &lib.TokenRes{}
	if err = poster.Handle(&lib.Token{Token: token}, tokenRes); err != nil { // 验证 token 请求错误
		return
	}

	if tokenRes.Code < 0 { // token 未验证成功
		return nil, errors.New("validate token error")
	}

	return
}

// home 页面是选择注册或者登录页面。如果本地存储中 token 验证合法后可自动登录并刷新本地 token，然后进入用户列表页面。如果本地没有 token，或者验证 token 失败，则进入 home 页面。
func renderHome(poster lib.Post, storage *storage_t) (renderHome bool) {
	tokenRes, err := validateToken(poster, storage)
	if err != nil {
		return true
	}

	if err = storage.StoreToken(tokenRes); err != nil {
		return true
	}

	return
}

// 渲染 UI
func renderUI(poster lib.Post, storage *storage_t, sigChan chan<- os.Signal) {
	b := initialBase(poster, storage)
	defer b.close()

	var m tea.Model
	if renderHome(poster, storage) {
		m = initialHome(b)
	} else {
		m = initialUsers(b)
	}

	p := tea.NewProgram(m)

	_, err := p.Run()
	lib.FatalNotNil(err)

	sigChan <- os.Interrupt
}

// 存储文件名字环境变量
func dbName() string {
	dbName, found := os.LookupEnv("DB_NAME")
	if !found {
		dbName = "client.db"
	}
	return dbName
}

// 存储文件路径
func dbPath() string {
	return filepath.Join(lib.WorkDir, dbName())
}

// 服务器地址环境变量
func svrAddr() string {
	svrAddr, found := os.LookupEnv("SVR_ADDR")
	if !found {
		svrAddr = "127.0.0.1:8888"
	}
	return svrAddr
}

// 连接服务器，连接失败按照指数回退策略重试，最多重试20次
func connect(sigChan <-chan os.Signal) (net.Conn, error) {
	for i := 0; i < 15; i++ {
		select {
		// 如果 UI 已退出，停止连接服务器
		case <-sigChan:
			return nil, nil
		default:
			// 客户端进行 tcp 拨号，请求连接服务器
			if conn, err := net.DialTimeout("tcp", svrAddr(), 3*time.Second); err == nil {
				return conn, nil
			}

			p := math.Pow(2, float64(i))
			r := float64(rand.Intn(1000))
			d := time.Duration(math.Min(p+r, 8000))

			// sleep
			time.Sleep(d * time.Millisecond)
		}
	}
	return nil, errors.New("connect error")
}

func main() {
	// 创建信号 channel
	sigChan := make(chan os.Signal, 1)
	// 注册要监听哪些信号
	signal.Notify(sigChan, os.Interrupt)    // ctrl+c
	signal.Notify(sigChan, syscall.SIGTERM) // kill

	// 初始化存储
	storage, err := new(storage_t).Init(dbPath())
	lib.FatalNotNil(err)

	// 请求 channel
	reqChan := make(chan *request_t, 1024)

	// 渲染 UI
	var poster lib.Post = newPoster(reqChan)
	go renderUI(poster, storage, sigChan)

	var reconnect bool

	for {
		// 连接服务器
		conn, err := connect(sigChan)
		lib.FatalNotNil(err)
		if conn == nil { // quit UI
			return
		}

		// 重新连接需要验证 token
		if reconnect {
			go func() {
				if tokenRes, err := validateToken(poster, storage); err == nil {
					storage.StoreToken(tokenRes)
				}
			}()
			reconnect = false
		}

		// 响应 channel
		resChan := make(chan *response_t, 1024)

		// 启动单独的协程，接收处理或转发来自服务器的 packet
		go recvFrom(conn, resChan, storage)

		// 当前协程调用并阻塞与 sendTo 函数，发送请求并接收响应
		if quit := sendTo(conn, reqChan, resChan); quit {
			return
		}

		reconnect = true
	}
}
