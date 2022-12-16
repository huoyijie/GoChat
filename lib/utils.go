package lib

import (
	"bufio"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func FatalNotNil(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func LogNotNil(err error) {
	if err != nil {
		log.Println(err)
	}
}

func LogMessage(msg ...any) {
	log.Println(msg...)
}

func LogConn() func() {
	log.Println("conn created")
	return func() {
		log.Println("conn closed")
	}
}

func HandleConnection(conn net.Conn) {
	defer conn.Close()
	defer LogConn()()

	scanner := bufio.NewScanner(conn)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		bytes := scanner.Bytes()
		if len(bytes) > 0 {
			LogMessage(string(bytes))
		}
	}

	LogNotNil(scanner.Err())
}

func SignalHandler() {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)
	signal.Notify(sigchan, syscall.SIGTERM)
	<-sigchan
	defer os.Exit(0)
}