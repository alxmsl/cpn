package strategies

import (
	"context"

	"github.com/alxmsl/cpn"
)

// PassFunc receives one token and returns one token
type PassFunc func(context.Context, *cpn.M) *cpn.M

type pass struct {
	chin  chan *cpn.M
	chout chan *cpn.M

	f PassFunc
}

// PassFuncOption creates a 1->1 strategy option for a place
func PassFuncOption(f PassFunc) cpn.StrategyOption {
	return passOption{f}
}

type passOption struct {
	f PassFunc
}

func (o passOption) Apply(p cpn.Strategy) {
	p.(*pass).f = o.f
}

// NewPass returns a 1->1 strategy. The 1->1 strategy receives one token and passes it forward. Token is modified in the
// PassFunc implementation
func NewPass(opts ...cpn.StrategyOption) cpn.Strategy {
	r := &pass{
		chin:  make(chan *cpn.M),
		chout: make(chan *cpn.M),
	}
	for _, o := range opts {
		o.Apply(r)
	}
	return r
}

func (p *pass) In() chan<- *cpn.M {
	return p.chin
}

func (p *pass) Out() <-chan *cpn.M {
	return p.chout
}

func (p *pass) Run() {
	defer close(p.chout)
	for m := range p.chin {
		p.chout <- p.f(context.Background(), m)
	}
}
