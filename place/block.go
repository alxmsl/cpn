package place

import "github.com/alxmsl/rtpn"

type Block struct {
	ch chan *rtpn.M
}

func NewBlock() *Block {
	return &Block{make(chan *rtpn.M)}
}

func (b Block) In() chan<- *rtpn.M {
	return b.ch
}

func (b Block) Out() <-chan *rtpn.M {
	return b.ch
}
