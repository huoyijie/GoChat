package lib

import (
	"github.com/bwmarrin/snowflake"
)

// 封装客户端连接，增加 snowflake.ID
type Socket struct {
	Id       snowflake.ID
	PackChan chan<- *Packet
}
