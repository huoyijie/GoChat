package main

import (
    "fmt"
    "net"

    "github.com/huoyijie/GoChat/lib"
)

func main() {
    // 客户端进行 tcp 拨号，请求连接 127.0.0.1:8888
    conn, err := net.Dial("tcp", "127.0.0.1:8888")
    // 连接遇到错误则退出进程
    lib.FatalNotNil(err)

    // 连接服务端成功，启动单独协程处理另一侧发送过来的消息
    go lib.HandleConnection(conn)

    // 向服务端发送消息
    fmt.Fprintf(conn, "Hello, World!\r\n")

    // 阻塞主线程，直到收到 ctrl+c 或者 kill 信号，退出进程
    lib.SignalHandler()
}
