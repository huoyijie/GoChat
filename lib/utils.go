package lib

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/snowflake"
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

// 输出日志
func LogMessage(msg ...any) {
	log.Println(msg...)
}

// 打印消息
func PrintMessage(msg *Msg) {
	fmt.Fprintf(os.Stdout, "%d->%d:%s\n", msg.From, msg.To, msg.Data)
}

// 输出连接建立与关闭消息到日志，包内私有方法，外部不能调用
func logConn(id snowflake.ID) func() {
	fmt.Fprintf(os.Stdout, "%d: 已连接\n", id)
	return func() {
		fmt.Fprintf(os.Stdout, "%d: 已断开连接\n", id)
	}
}

// 接收连接另一侧发送的消息，输出消息到日志
func HandleConnection(conn net.Conn, id snowflake.ID, handleMsg func(*Msg), close func()) {
	// 连接建立和断开时，分别输出日志
	defer logConn(id)()

	// 从当前方法返回后，断开连接，清理资源等
	defer close()

	// 设置如何处理接收到的字节流，SplitFunc 会根据 packet 开头 length 把字节流分割为消息流
	scanner := bufio.NewScanner(conn)
	scanner.Split(SplitFunc)

	// 循环解析消息，每当解析出一条消息后，scan() 返回 true
	for scanner.Scan() {
		// 把 scanner 解析出的消息字节 slice 转换为 Msg
		msg, err := RecvMsg(scanner.Bytes())
		if err != nil {
			return
		}
		handleMsg(msg)
	}

	// 如果解析消息遇到错误，则输出错误到日志
	LogNotNil(scanner.Err())
}

// 信号监听处理器
func SignalHandler() {
	// 创建信号 channel
	sigchan := make(chan os.Signal, 1)

	// 注册要监听哪些信号
	signal.Notify(sigchan, os.Interrupt)    // ctrl+c
	signal.Notify(sigchan, syscall.SIGTERM) // kill

	// 一直阻塞，直到收到信号，恢复执行并退出进程
	<-sigchan
	// 退出进程
	defer os.Exit(0)
}
