package lib

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var HomeDir = getHomeDir()
var WorkDir = getWorkDir()

// 返回 home 目录 ~/
func getHomeDir() (homeDir string) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	return
}

// 返回工作目录 ~/.gochat
func getWorkDir() (workDir string) {
	workDir = filepath.Join(HomeDir, ".gochat")
	if _, err := os.Stat(workDir); os.IsNotExist(err) {
		if err := os.Mkdir(workDir, 00744); err != nil {
			log.Fatal(err)
		}
	}
	return
}

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
	fmt.Fprintf(os.Stdout, "%s->%s:%s\n", msg.From, msg.To, msg.Data)
}

// 打印服务器错误
func PrintErr(errRes *ErrRes) {
	fmt.Fprintf(os.Stdout, "系统异常: %d\n", errRes.Code)
}
