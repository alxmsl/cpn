package strategies

import (
	"context"
	"fmt"

	ctx_pkg "github.com/alxmsl/cpn/context"

	"github.com/alxmsl/cpn"
)

// ForkFunc receives one token and returns many tokens
type ForkFunc func(context.Context, *cpn.M, chan<- *cpn.M)

type fork struct {
	chin  chan *cpn.M
	chout chan *cpn.M

	f ForkFunc
}

// ForkFuncOption creates a 1->m strategy option for a place
func ForkFuncOption(f ForkFunc) cpn.StrategyOption {
	return forkOption{f}
}

type forkOption struct {
	f ForkFunc
}

func (o forkOption) Apply(p cpn.Strategy) {
	p.(*fork).f = o.f
}

// NewFork returns a 1->m strategy. The 1->m strategy receives one token and passes many tokens forward. Token is
// modified in the ForkFunc implementation
func NewFork(opts ...cpn.StrategyOption) cpn.Strategy {
	r := &fork{
		chin:  make(chan *cpn.M),
		chout: make(chan *cpn.M),
	}
	for _, o := range opts {
		o.Apply(r)
	}
	return r
}

func (p *fork) In() chan<- *cpn.M {
	return p.chin
}

func (p *fork) Out() <-chan *cpn.M {
	return p.chout
}

func (p *fork) Run(ctx context.Context) {
	defer close(p.chout)
	dCtx := ctx_pkg.Detach(ctx)
	skip := false
	for m := range p.chin {
		select {
		case <-ctx.Done():
			fmt.Println(">>> DONE")
			skip = true
		default:
		}
		if !skip {
			p.f(dCtx, m, p.chout)
		}
	}
}
