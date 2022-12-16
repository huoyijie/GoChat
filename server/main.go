package main

import (
	"fmt"
	"net"

	"github.com/huoyijie/GoChat/lib"
)

func main() {
    addr := ":8888"
    go lib.SignalHandler()

	ln, err := net.Listen("tcp", addr)
	lib.FatalNotNil(err)
    lib.LogMessage("Listening on", addr)

	for {
		conn, err := ln.Accept()
		lib.FatalNotNil(err)

        fmt.Fprintf(conn, "http/1.0 200 ok\r\n\r\n")
		go lib.HandleConnection(conn)
	}
}
