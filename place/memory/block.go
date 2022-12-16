package memory

import (
	"context"

	"github.com/alxmsl/cpn"
)

type Block struct {
	ch chan *cpn.M
}

func NewBlock(_ ...cpn.StrategyOption) cpn.Strategy {
	return &Block{make(chan *cpn.M)}
}

func (b Block) Run(_ context.Context) {}

func (b Block) In() chan<- *cpn.M {
	return b.ch
}

func (b Block) Out() <-chan *cpn.M {
	return b.ch
}
