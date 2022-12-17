package main

import (
    "fmt"
    "net"

    "github.com/huoyijie/GoChat/lib"
)

func main() {
    // 启动单独协程，监听 ctrl+c 或 kill 信号，收到信号结束进程
    go lib.SignalHandler()

    // tcp 监听地址 0.0.0.0:8888
    addr := ":8888"

    // tcp 监听
    ln, err := net.Listen("tcp", addr)

    // tcp 监听遇到错误退出进程
    lib.FatalNotNil(err)
    // 输出日志
    lib.LogMessage("Listening on", addr)

    // 循环接受客户端连接
    for {
        // 每当有客户端连接时，ln.Accept 会返回新的连接 conn
        conn, err := ln.Accept()
        // 如果接受的新连接遇到错误，则退出进程
        lib.FatalNotNil(err)

        // 通过 conn 向连接的另一侧发送消息
        fmt.Fprintf(conn, "Hello, World!\r\n")

        // 为每个新连接启动一个单独协程，该协程会读取另一侧发送的消息
        go lib.HandleConnection(conn)
    }
}
