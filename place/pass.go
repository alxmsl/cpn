package place

import (
	"context"

	"github.com/alxmsl/cpn"
)

type (
	PassFunc func(context.Context, *cpn.M)

	pass struct {
		chin  chan *cpn.M
		chout chan *cpn.M

		f PassFunc
	}
)

// PassFuncOption returns a place option contains a PassFunc
func PassFuncOption(f PassFunc) cpn.PlaceOption {
	return passOption{f}
}

type passOption struct {
	f PassFunc
}

func (o passOption) Apply(p cpn.Place) {
	p.(*pass).f = o.f
}

// NewPass returns a Pass Place. The Pass Place gets one mark and passes it forward. It is allowed to modify a mark's
// value using function PassFunc
func NewPass(opts ...cpn.PlaceOption) cpn.Place {
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
		p.f(context.Background(), m)
		p.chout <- m
	}
}
