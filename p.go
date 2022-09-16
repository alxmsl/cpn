package cpn

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/alxmsl/prmtvs/skm"
)

// Strategy allows to create special place strategies.
// Strategy is defined as two channels: in and out; and run handler (@todo: maybe stop?)
type Strategy interface {
	In() chan<- *M
	Out() <-chan *M
	Run()
}

// P implements an abstract place in PN
type P struct {
	ctx context.Context
	mu  sync.Mutex

	// name is a place name in the PN. Should be unique
	name string

	// strategy defines a behaviour for the place. Basic behaviours are (@todo: link places here):
	//  - 1->1 means one token ins, then one token outs
	//  - 1->m means once token ins, many tokens out
	//  - m->1 means many tokens in, one token outs
	// But this is allowed to implement own Strategy
	strategy Strategy

	// ins is a sorted set of incoming edges
	ins *skm.SKM
	// out is a channel for outgoing edges
	out chan *M

	// State flags for an abstract place:
	//  - i means an initial place in the PN. Initial place doesn't have incoming edges
	//  - t means a terminal place in the PN. Terminal place doesn't have outgoing edges
	//  - k means to don't clean up Strategy object in the terminal place. This is used for test purposes
	i, t, k bool

	// s keeps an internal state for the place:
	//  - stateActive means place started receiving tokens (@todo: probably, we don't need this state)
	//  - stateClosed means place stopped sending tokens. Reasons: Strategy object stopped working or the place
	//    context deadline exceeded
	//  - stateReady means place is ready to accept tokens
	s int64
}

// NewP creates a new place with specific name. By default, place is both initial and terminal. When place is linked to
// transition default state is changed
func NewP(name string) *P {
	return &P{
		name: name,

		ins: skm.NewSKM(),
		out: make(chan *M),

		i: true,
		t: true,
	}
}

func (p *P) SetOptions(opts ...PlaceOption) *P {
	for _, opt := range opts {
		opt.Apply(p)
	}
	return p
}

func (p *P) Name() string {
	return p.name
}

func (p *P) Close() {
	close(p.strategy.In())
}

func (p *P) In() chan<- *M {
	return p.strategy.In()
}

func (p *P) Out() <-chan *M {
	if !p.t {
		return nil
	}
	return p.strategy.Out()
}

func (p *P) ready() bool {
	s := atomic.LoadInt64(&p.s)
	return s&(stateReady|stateClosed) > 0x0
}

func (p *P) run() {
	p.strategy.Run()
}

func (p *P) recv() {
	atomic.AddInt64(&p.s, stateActive)
	if p.i {
		return
	}

	wg := &sync.WaitGroup{}
	wg.Add(p.ins.Len())
	p.ins.Over(func(i int, n string, v interface{}) bool {
		go func() {
			defer wg.Done()
			for m := range v.(chan *M) {
				m.passP(p)
				p.strategy.In() <- m
			}
		}()
		return true
	})
	wg.Wait()

	close(p.strategy.In())
	if p.t {
		close(p.out)
	}
}

func (p *P) send() {
	if p.t {
		if !p.k {
			for m := range p.strategy.Out() {
				m.passP(p)
			}
		}
		return
	}
	for atomic.LoadInt64(&p.s)&stateClosed == 0x0 {
		select {
		case m, ok := <-p.strategy.Out():
			if !ok {
				atomic.AddInt64(&p.s, stateClosed)
				break
			}
			atomic.AddInt64(&p.s, stateReady)

			m.passP(p)
			p.out <- m

			atomic.AddInt64(&p.s, -stateReady)
			p.mu.Unlock()
		case <-p.ctx.Done():
			atomic.AddInt64(&p.s, stateClosed)
		}
	}
	close(p.out)
}

const (
	stateActive = 1 << 0
	stateClosed = 1 << 1
	stateReady  = 1 << 2
)
