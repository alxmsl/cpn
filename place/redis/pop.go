package redis

import (
	"reflect"

	"github.com/alxmsl/cpn"
	"github.com/mediocregopher/radix/v3"
)

type Pop struct {
	chin  chan *cpn.M
	chout chan *cpn.M

	errs chan<- error
	f    UnmarshalFunc
	key  string
	pool *radix.Pool
	t    reflect.Type
}

func NewPop(opts ...cpn.StrategyOption) cpn.Strategy {
	p := &Pop{
		chin:  make(chan *cpn.M),
		chout: make(chan *cpn.M),
	}
	for _, o := range opts {
		o.Apply(p)
	}
	return p
}

func (p *Pop) In() chan<- *cpn.M {
	return p.chin
}

func (p *Pop) Out() <-chan *cpn.M {
	return p.chout
}

func (p *Pop) SetErrs(errs chan<- error) {
	p.errs = errs
}

func (p *Pop) SetKey(k string) {
	p.key = k
}

func (p *Pop) SetPool(pool *radix.Pool) {
	p.pool = pool
}

func (p *Pop) SetType(t reflect.Type) {
	p.t = t
}

func (p *Pop) Run() {
	defer close(p.chout)
	var (
		err error
		s   string
	)
	for m := range p.chin {
		v := reflect.New(p.t).Interface()
		err = p.pool.Do(radix.Cmd(&s, "RPOP", p.key))
		if err == nil {
			err = p.f(s, &v)
		}
		if err != nil {
			select {
			case p.errs <- err:
			default:
			}
			continue
		}
		m.SetValue(v)
		p.chout <- m
	}
}
