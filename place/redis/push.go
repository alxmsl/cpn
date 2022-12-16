package redis

import (
	"context"

	"github.com/alxmsl/cpn"
	"github.com/mediocregopher/radix/v3"
)

type Push struct {
	chin  chan *cpn.M
	chout chan *cpn.M

	errs chan<- error
	f    MarshalFunc
	key  string
	pool *radix.Pool
}

func NewPush(opts ...cpn.StrategyOption) cpn.Strategy {
	p := &Push{
		chin:  make(chan *cpn.M),
		chout: make(chan *cpn.M),
	}
	for _, o := range opts {
		o.Apply(p)
	}
	return p
}

func (p *Push) In() chan<- *cpn.M {
	return p.chin
}

func (p *Push) Out() <-chan *cpn.M {
	return p.chout
}

func (p *Push) SetErrs(errs chan<- error) {
	p.errs = errs
}

func (p *Push) SetKey(k string) {
	p.key = k
}

func (p *Push) SetPool(pool *radix.Pool) {
	p.pool = pool
}

func (p *Push) Run(_ context.Context) {
	defer close(p.chout)
	for m := range p.chin {
		v, err := p.f(m.Value())
		if err == nil {
			err = p.pool.Do(radix.Cmd(nil, "LPUSH", p.key, v))
		}
		if err != nil {
			select {
			case p.errs <- err:
			default:
			}
			continue
		}
		p.chout <- m
	}
}
