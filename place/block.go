package place

import "github.com/alxmsl/rtpn"

type Block struct {
	ch chan *rtpn.M
}

func NewBlock() *Block {
	return &Block{make(chan *rtpn.M)}
}

func (b Block) ReadCh() <-chan *rtpn.M {
	return b.ch
}

func (b Block) WriteCh() chan<- *rtpn.M {
	return b.ch
}
