package main

import (
	"fmt"
	"net"

	"github.com/bwmarrin/snowflake"
	"github.com/huoyijie/GoChat/lib"
)

func main() {
	// 启动单独协程，监听 ctrl+c 或 kill 信号，收到信号结束进程
	go lib.SignalHandler()

	// 客户端进行 tcp 拨号，请求连接 127.0.0.1:8888
	conn, err := net.Dial("tcp", "127.0.0.1:8888")
	// 连接遇到错误则退出进程
	lib.FatalNotNil(err)

	// id 由服务器生成，暂时未发给客户端
	var id snowflake.ID
	// 连接成功后启动协程输出服务器的转发消息
	go lib.HandleConnection(
		conn,
		id,
		func(msg string) {
			lib.PrintMessage(msg)
		},
		func() {
			// 从当前方法返回时，关闭连接
			conn.Close()
		})

	var input string
	for {
		// 读取用户输入消息
		fmt.Scanf("%s", &input)
		// 向服务端发送消息
		fmt.Fprintf(conn, "%s\r\n", input)
	}
}
