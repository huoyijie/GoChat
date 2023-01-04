package main

import (
	"net"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/huoyijie/GoChat/lib"
)

func main() {
	// 启动单独协程，监听 ctrl+c 或 kill 信号，收到信号结束进程
	go lib.SignalHandler()

	// 客户端进行 tcp 拨号，请求连接 127.0.0.1:8888
	conn, err := net.Dial("tcp", "127.0.0.1:8888")
	// 连接遇到错误则退出进程
	lib.FatalNotNil(err)

	packChan := make(chan *lib.Packet)
	go lib.SendTo(conn, packChan)

	// 连接成功后启动协程输出服务器的转发消息
	go lib.RecvFrom(
		conn,
		func(pack *lib.Packet) {
			if pack.Kind == lib.PackKind_MSG {
				msg := &lib.Msg{}
				if err := lib.Unmarshal(pack.Data, msg); err == nil {
					lib.PrintMessage(msg)
				}
			}
		},
		func() {
			// 从当前方法返回时，关闭连接
			conn.Close()
		})

	p := tea.NewProgram(home{choice: CHOICE_SIGNIN, base: base{packChan: packChan}})

	_, err = p.Run()
	lib.FatalNotNil(err)
}
