package main

import (
	"net"
	"path/filepath"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/huoyijie/GoChat/lib"
)

type request_t struct {
	pack     *lib.Packet
	c        chan *response_t
	deadline time.Time
}

func (request *request_t) init(pack *lib.Packet, sync bool) *request_t {
	request.pack = pack
	if sync {
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
			for id, request := range requests {
				if time.Now().After(request.deadline) {
					request.c <- new(response_t)
					delete(requests, id)
				}
			}
		}
	}
}

func main() {
	// 启动单独协程，监听 ctrl+c 或 kill 信号，收到信号结束进程
	go lib.SignalHandler()

	storage, err := new(Storage).Init(filepath.Join(lib.WorkDir, "client.db"))
	lib.FatalNotNil(err)

	// 客户端进行 tcp 拨号，请求连接 127.0.0.1:8888
	conn, err := net.Dial("tcp", "127.0.0.1:8888")
	// 连接遇到错误则退出进程
	lib.FatalNotNil(err)

	reqChan := make(chan *request_t)
	resChan := make(chan *response_t)
	go sendTo(conn, reqChan, resChan)

	// 连接成功后启动协程输出服务器的转发消息
	go lib.RecvFrom(
		conn,
		func(pack *lib.Packet) error {
			switch pack.Kind {
			case lib.PackKind_MSG:
				msg := &lib.Msg{}
				if err := lib.Unmarshal(pack.Data, msg); err == nil {
					lib.PrintMessage(msg)
				}
			default:
				resChan <- &response_t{pack: pack}
			}
			return nil
		},
		func() {
			// 从当前方法返回时，关闭连接
			conn.Close()
		})

	p := tea.NewProgram(home{choice: CHOICE_SIGNIN, base: base{reqChan: reqChan, storage: storage}})

	_, err = p.Run()
	lib.FatalNotNil(err)
}
