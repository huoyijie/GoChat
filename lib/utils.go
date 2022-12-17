package lib

import (
    "bufio"
    "log"
    "net"
    "os"
    "os/signal"
    "syscall"
)

// 如果 err != nil，输出错误日志并退出进程
func FatalNotNil(err error) {
    if err != nil {
        log.Fatal(err)
    }
}

// 如果 err != nil，输出错误日志
func LogNotNil(err error) {
    if err != nil {
        log.Println(err)
    }
}

// 输出消息到日志
func LogMessage(msg ...any) {
    log.Println(msg...)
}

// 输出连接建立与关闭消息到日志，包内私有方法，外部不能调用
func logConn() func() {
    log.Println("已连接")
    return func() {
        log.Println("已断开连接")
    }
}

// 接收连接另一侧发送的消息，输出消息到日志
func HandleConnection(conn net.Conn) {
    // 连接建立和断开时，分别输出日志
    defer logConn()()
    // 从当前方法返回时，关闭连接
    defer conn.Close()

    // 设置如何处理接收到的字节流，bufio.ScanLines 为逐行扫描的方式把字节流分割为消息流
    scanner := bufio.NewScanner(conn)
    scanner.Split(bufio.ScanLines)

    // 循环解析消息，每当解析出一条消息后，scan() 返回 true
    for scanner.Scan() {
        // 返回解析出的消息字节 slice
        bytes := scanner.Bytes()
        // 消息内容不为空，则输出到日志
        if len(bytes) > 0 {
            LogMessage(string(bytes))
        }
    }

    // 如果解析消息遇到错误，则输出错误到日志
    LogNotNil(scanner.Err())
}

// 信号监听处理器
func SignalHandler() {
    // 创建信号 channel
    sigchan := make(chan os.Signal, 1)

    // 注册要监听哪些信号
    signal.Notify(sigchan, os.Interrupt) // ctrl+c
    signal.Notify(sigchan, syscall.SIGTERM) // kill

    // 一直阻塞，直到收到信号，恢复执行并退出进程
    <-sigchan
    // 退出进程
    defer os.Exit(0)
}