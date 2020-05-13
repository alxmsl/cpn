package place

import "github.com/alxmsl/cpn"

type Block struct {
	ch chan *cpn.M
}

func NewBlock() cpn.Place {
	return &Block{make(chan *cpn.M)}
}

func (b Block) Close() {
	close(b.ch)
}

func (b Block) Run() {}

func (b Block) In() chan<- *cpn.M {
	return b.ch
}

func (b Block) Out() <-chan *cpn.M {
	return b.ch
}
