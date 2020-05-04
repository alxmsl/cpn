package place

import "github.com/alxmsl/pn"

type Block struct {
	ch chan *pn.M
}

func NewBlock() *Block {
	return &Block{make(chan *pn.M)}
}

func (b Block) ReadCh() <-chan *pn.M {
	return b.ch
}

func (b Block) WriteCh() chan<- *pn.M {
	return b.ch
}
