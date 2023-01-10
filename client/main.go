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

func sendTo(conn net.Conn, reqChan <-chan *request_t, resChan <-chan *response_t) {
	var id uint64
	requests := make(map[uint64]*request_t)
	timeoutTicker := time.NewTicker(50 * time.Millisecond)
	defer timeoutTicker.Stop()

	for {
		select {
		case req := <-reqChan:
			id++
			req.pack.Id = id
			bytes, err := lib.MarshalPack(req.pack)
			if err == nil {
				_, err = conn.Write(bytes)
			}

			if err != nil {
				return
			}

			if req.sync() {
				req.deadline = time.Now().Add(5 * time.Second)
				requests[req.pack.Id] = req
			}
		case res := <-resChan:
			if req, ok := requests[res.pack.Id]; ok {
				req.c <- res
				delete(requests, res.pack.Id)
			}
		case <-timeoutTicker.C:
			for id, req := range requests {
				if time.Now().After(req.deadline) {
					req.c <- new(response_t)
					delete(requests, id)
				}
			}
		}
	}
}

func recvFrom(conn net.Conn, msgChan chan<- *lib.Msg, resChan chan<- *response_t) {
	lib.RecvFrom(
		conn,
		func(pack *lib.Packet) error {
			switch pack.Kind {
			case lib.PackKind_MSG:
				msg := &lib.Msg{}
				if err := lib.Unmarshal(pack.Data, msg); err == nil {
					go func() {
						msgChan <- msg
					}()
				}
			case lib.PackKind_ERR:
				errRes := &lib.ErrRes{}
				if err := lib.Unmarshal(pack.Data, errRes); err == nil {
					lib.FatalNotNil(fmt.Errorf("系统异常: %d", errRes.Code))
				}
			case lib.PackKind_RES:
				resChan <- &response_t{pack: pack}
			}
			return nil
		})
}

// home 页面是选择注册或者登录页面。如果本地存储中 token 验证合法后可自动登录并刷新本地 token，然后进入用户列表页面。如果本地没有 token，或者验证 token 失败，则进入 home 页面。
func renderHome(reqChan chan<- *request_t, storage *Storage) (renderHome bool) {
	kv, err := storage.GetValue("token")
	if err != nil { // 未登录过
		return true
	}

	token, err := base64.StdEncoding.DecodeString(kv.Value)
	if err != nil { // token 解析错误
		return true
	}

	bytes, err := lib.Marshal(&lib.Token{Token: token})
	if err != nil { // 序列化 token 错误
		return true
	}

	req := newRequest(&lib.Packet{Kind: lib.PackKind_TOKEN, Data: bytes})
	reqChan <- req
	res := <-req.c
	if !res.ok() { // 验证 token 请求超时
		return true
	}

	tokenRes := &lib.TokenRes{}
	err = lib.Unmarshal(res.pack.Data, tokenRes)
	if err != nil || tokenRes.Code < 0 { // 验证 token 响应错误或者 token 未验证成功
		return true
	}

	err = storage.NewKVS([]KeyValue{
		{Key: "id", Value: fmt.Sprintf("%d", tokenRes.Id)},
		{Key: "username", Value: tokenRes.Username},
		{Key: "token", Value: base64.StdEncoding.EncodeToString(tokenRes.Token)},
	})
	if err != nil { // token 验证成功，但是写入存储错误
		return true
	}

	return
}

// 渲染 UI
func renderUI(reqChan chan<- *request_t, msgChan <-chan *lib.Msg, storage *Storage) {
	b := base{reqChan: reqChan, msgChan: msgChan, storage: storage}

	var m tea.Model
	if renderHome(reqChan, storage) {
		m = home{choice: CHOICE_SIGNIN, base: b}
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
	storage, err := new(Storage).Init(filepath.Join(lib.WorkDir, dbName()))
	lib.FatalNotNil(err)

	// 客户端进行 tcp 拨号，请求连接服务器
	conn, err := net.Dial("tcp", svrAddr())
	// 连接遇到错误则退出进程
	lib.FatalNotNil(err)

	reqChan := make(chan *request_t)
	resChan := make(chan *response_t)
	msgChan := make(chan *lib.Msg)

	// 启动单独的协程，发送请求并接收响应
	go sendTo(conn, reqChan, resChan)

	// 启动单独的协程，接收处理或转发来自服务器的 packet
	go recvFrom(conn, msgChan, resChan)

	// 渲染 UI
	renderUI(reqChan, msgChan, storage)
}
