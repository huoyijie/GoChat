package main

import (
	"encoding/base64"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/huoyijie/GoChat/lib"
)

// 向服务器发送 packet。如果是同步请求，会通过 request.c 返回服务器响应数据，同时也会检查同步请求是否已超时。
func sendTo(conn net.Conn, reqChan <-chan *request_t, resChan <-chan *response_t) {
	var id uint64
	// 登记所有的同步请求，并等待响应
	requests := make(map[uint64]*request_t)
	// 同步请求超时检查间隔 50ms
	timeoutTicker := time.NewTicker(50 * time.Millisecond)
	defer timeoutTicker.Stop()

	for {
		select {
		// 有服务器请求进来
		case request := <-reqChan:
			// 分配 packet.id
			id++
			request.pack.Id = id

			bytes, err := lib.MarshalPack(request.pack)
			if err != nil { // 序列化 packet 错误
				return
			}

			if _, err = conn.Write(bytes); err != nil { // 发送字节数据错误
				return
			}

			if request.sync() { // 登记同步请求
				// 如果 5s 内服务器没有返回，则超时
				request.deadline = time.Now().Add(5 * time.Second)
				requests[request.pack.Id] = request
			}
		// 通过 resChan 接收从 recvFrom 协程发送过来的响应对象
		case response := <-resChan:
			if request, ok := requests[response.pack.Id]; ok {
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
		}
	}
}

// 从服务器接收 packet 并进行处理
func recvFrom(conn net.Conn, resChan chan<- *response_t, storage *Storage) {
	lib.RecvFrom(
		conn,
		// 从服务器接收 packet 的处理函数
		func(pack *lib.Packet) error {
			switch pack.Kind {
			// 当前连接的登录用户收到新未读消息
			case lib.PackKind_MSG:
				msg := &lib.Msg{}
				if err := lib.Unmarshal(pack.Data, msg); err == nil {
					// 新消息写入本地存储
					storage.NewMsg(&Message{
						Id:   msg.Id,
						Kind: int32(msg.Kind),
						From: msg.From,
						Data: msg.Data,
					})
				}
			// 当前连接遇到系统异常，退出进程
			case lib.PackKind_ERR:
				errRes := &lib.ErrRes{}
				if err := lib.Unmarshal(pack.Data, errRes); err == nil {
					lib.FatalNotNil(fmt.Errorf("系统异常: %d", errRes.Code))
				}
			// 收到同步请求的响应
			case lib.PackKind_RES:
				// 当前是在 recvFrom 协程里，需要把 packet 封装为 response_t 对象，并通过 resChan channel 发到 sendTo 协程
				resChan <- newResponse(pack)
			}
			return nil
		})
}

// home 页面是选择注册或者登录页面。如果本地存储中 token 验证合法后可自动登录并刷新本地 token，然后进入用户列表页面。如果本地没有 token，或者验证 token 失败，则进入 home 页面。
func renderHome(poster lib.Post, storage *Storage) (renderHome bool) {
	kv, err := storage.GetValue("token")
	if err != nil { // 未登录过
		return true
	}

	token, err := base64.StdEncoding.DecodeString(kv.Value)
	if err != nil { // token 解析错误
		return true
	}

	tokenRes := &lib.TokenRes{}
	if err = poster.Handle(&lib.Token{Token: token}, tokenRes); err != nil || tokenRes.Code < 0 { // 验证 token 请求错误，或者 token 未验证成功
		return true
	}

	if err = storage.NewKVS([]KeyValue{
		{Key: "id", Value: fmt.Sprintf("%d", tokenRes.Id)},
		{Key: "username", Value: tokenRes.Username},
		{Key: "token", Value: base64.StdEncoding.EncodeToString(tokenRes.Token)},
	}); err != nil { // token 验证成功，但是写入存储错误
		return true
	}

	return
}

// 渲染 UI
func renderUI(poster lib.Post, storage *Storage) {
	b := initialBase(poster, storage)

	var m tea.Model
	if renderHome(poster, storage) {
		m = initialHome(b)
	} else {
		m = initialUsers(b)
	}

	p := tea.NewProgram(m)

	_, err := p.Run()
	lib.FatalNotNil(err)
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

// 删除本地存储文件
func dropDB() error {
	return os.Remove(dbPath())
}

// 服务器地址环境变量
func svrAddr() string {
	svrAddr, found := os.LookupEnv("SVR_ADDR")
	if !found {
		svrAddr = "127.0.0.1:8888"
	}
	return svrAddr
}

func main() {
	// 启动单独协程，监听 ctrl+c 或 kill 信号，收到信号结束进程
	go lib.SignalHandler()

	// 初始化存储
	storage, err := new(Storage).Init(dbPath())
	lib.FatalNotNil(err)

	// 客户端进行 tcp 拨号，请求连接服务器
	conn, err := net.Dial("tcp", svrAddr())
	// 连接遇到错误则退出进程
	lib.FatalNotNil(err)

	reqChan := make(chan *request_t)
	resChan := make(chan *response_t)

	// 启动单独的协程，发送请求并接收响应
	go sendTo(conn, reqChan, resChan)

	// 启动单独的协程，接收处理或转发来自服务器的 packet
	go recvFrom(conn, resChan, storage)

	// 渲染 UI
	renderUI(newPoster(reqChan), storage)
}
