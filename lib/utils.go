package lib

import (
	"bufio"
	"fmt"
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

// 输出日志
func LogMessage(msg ...any) {
	log.Println(msg...)
}

// 打印消息
func PrintMessage(msg *Msg) {
	fmt.Fprintf(os.Stdout, "%d->%d:%s\n", msg.From, msg.To, msg.Data)
}

// 接收连接另一侧发送的消息，输出消息到日志
func RecvFrom(conn net.Conn, handlePack func(*Packet), close func()) {
	// 从当前方法返回后，断开连接，清理资源等
	defer close()

	// 设置如何处理接收到的字节流，SplitFunc 会根据 packet 开头 length 把字节流分割为消息流
	scanner := bufio.NewScanner(conn)
	scanner.Split(SplitFunc)

	// 循环解析消息，每当解析出一条消息后，scan() 返回 true
	for scanner.Scan() {
		// 把 scanner 解析出的消息字节 slice 解析为 Pack
		pack := &Packet{}
		if err := Unmarshal(scanner.Bytes(), pack); err != nil {
			return
		}
		handlePack(pack)
	}

	// 如果解析消息遇到错误，则输出错误到日志
	LogNotNil(scanner.Err())
}

func SendTo(conn net.Conn, ch <-chan *Packet) {
	var id uint64
	for pack := range ch {
		id++
		pack.Id = id

		bytes, err := MarshalPack(pack)
		if err == nil {
			_, err = conn.Write(bytes)
		}

		if err != nil {
			return
		}
	}
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
