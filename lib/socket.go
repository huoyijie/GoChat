package lib

type Socket struct {
	PackChan chan<- *Packet
}
