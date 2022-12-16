package http

import (
	"context"
	"net/http"

	"github.com/alxmsl/cpn"
)

func AddressOption(addr string) cpn.StrategyOption {
	return addressOption{addr}
}

type addressOption struct {
	addr string
}

func (o addressOption) Apply(p cpn.Strategy) {
	p.(*Request).addr = o.addr
}

func PatternOption(pattern string) cpn.StrategyOption {
	return patternOption{pattern}
}

type patternOption struct {
	pattern string
}

func (o patternOption) Apply(p cpn.Strategy) {
	p.(*Request).pattern = o.pattern
}

type Request struct {
	chin  chan *cpn.M
	chout chan *cpn.M

	addr, pattern string
	cancel        context.CancelFunc
}

func NewRequest(opts ...cpn.StrategyOption) cpn.Strategy {
	p := &Request{
		chin:  make(chan *cpn.M),
		chout: make(chan *cpn.M),
	}
	for _, o := range opts {
		o.Apply(p)
	}
	return p
}

func (p *Request) In() chan<- *cpn.M {
	return p.chin
}

func (p *Request) Out() <-chan *cpn.M {
	return p.chout
}

func (p *Request) SetCancel(cancel context.CancelFunc) {
	p.cancel = cancel
}

func (p *Request) Run(_ context.Context) {
	defer close(p.chout)
	http.HandleFunc(p.pattern, func(w http.ResponseWriter, r *http.Request) {
		ctx := &RequestContext{
			done: make(chan struct{}),
			r:    r,
			w:    w,
		}
		p.chout <- cpn.NewM(ctx)
		ctx.Wait()
	})
	if err := http.ListenAndServe(p.addr, nil); err != http.ErrServerClosed {
		p.cancel()
		return
	}
}
