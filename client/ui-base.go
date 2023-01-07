package main

import (
	"github.com/huoyijie/GoChat/lib"
)

type base struct {
	msgChan <-chan *lib.Msg
	reqChan chan<- *request_t
	storage *Storage
}
