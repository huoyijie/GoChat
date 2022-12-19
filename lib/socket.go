package lib

import (
	"net"

	"github.com/bwmarrin/snowflake"
)

// 封装客户端连接，增加 snowflake.ID
type Socket struct {
	Id   snowflake.ID
	Conn net.Conn
}