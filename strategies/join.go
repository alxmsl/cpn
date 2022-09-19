package strategies

import (
	"context"

	"github.com/alxmsl/cpn"
)

// JoinFunc receives many tokens and return just one token
type JoinFunc func(context.Context, <-chan *cpn.M) *cpn.M

type join struct {
	chin  chan *cpn.M
	chout chan *cpn.M

	f JoinFunc
}

// JoinFuncOption creates a m->1 strategy option for a place
func JoinFuncOption(f JoinFunc) cpn.StrategyOption {
	return joinOption{f}
}

type joinOption struct {
	f JoinFunc
}

func (o joinOption) Apply(p cpn.Strategy) {
	p.(*join).f = o.f
}

// NewJoin returns a m->1 strategy. The m->1 strategy receives many tokens and passes just one token forward. Token is
// modified in the JoinFunc implementation
func NewJoin(opts ...cpn.StrategyOption) cpn.Strategy {
	r := &join{
		chin:  make(chan *cpn.M),
		chout: make(chan *cpn.M),
	}
	for _, o := range opts {
		o.Apply(r)
	}
	return r
}

func (p *join) In() chan<- *cpn.M {
	return p.chin
}

func (p *join) Out() <-chan *cpn.M {
	return p.chout
}

func (p *join) Run() {
	defer close(p.chout)
	p.chout <- p.f(context.Background(), p.chin)
}
