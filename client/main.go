package main

import (
	"fmt"
	"net"

	"github.com/huoyijie/GoChat/lib"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8888")
	lib.FatalNotNil(err)

	go lib.HandleConnection(conn)

	fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")

	lib.SignalHandler()
}
