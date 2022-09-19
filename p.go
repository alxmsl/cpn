package cpn

import (
	"context"
	"github.com/alxmsl/cpn/trace"
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
	//  - k means to don't clean up Strategy object in the terminal place. This is used for test purposes
	//  - l means need to log an entity behaviour
	//  - t means a terminal place in the PN. Terminal place doesn't have outgoing edges
	i, k, l, t bool

	// s keeps an internal state for the place. See state constants for details
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
		l: trace.NeedLog(name),
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

func (p *P) Send(m *M) {
	if p.l {
		trace.Log(p.name, "[recv (direct)]", "v:", m.Value())
	}
	atomic.AddInt64(&p.s, stateProcessing)
	p.In() <- m
}

func (p *P) ready() bool {
	s := atomic.LoadInt64(&p.s)
	return s&(stateReady|stateClosed) > 0x0
}

func (p *P) run() {
	if p.l {
		trace.Log(p.name, "[running...]", "i:", p.i, "k:", p.k, "t:", p.t)
		defer trace.Log(p.name, "[running completed]")
	}
	p.strategy.Run()
}

func (p *P) recv() {
	if p.l {
		trace.Log(p.name, "[receiving...]")
		defer trace.Log(p.name, "[receiving completed]")
	}

	if p.i {
		if p.l {
			trace.Log(p.name, "[initial place detected]")
		}
		return
	}

	wg := &sync.WaitGroup{}
	wg.Add(p.ins.Len())
	p.ins.Over(func(i int, n string, v interface{}) bool {
		if p.l {
			trace.Log(p.name, "[listening...]", "n:", n)
		}
		go func() {
			if p.l {
				defer trace.Log(p.name, "[stop listening]", "n:", n)
			}
			defer wg.Done()
			for m := range v.(chan *M) {
				m.passP(p)
				if p.l {
					trace.Log(p.name, "[recv]", "n:", n, "v:", m.Value())
				}
				atomic.AddInt64(&p.s, stateProcessing)
				p.In() <- m
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
	if p.l {
		trace.Log(p.name, "[sending...]")
		defer trace.Log(p.name, "[sending completed]")
	}
	if p.t {
		if p.l {
			trace.Log(p.name, "[terminal place detected]")
		}
		if !p.k {
			trace.Log(p.name, "[utilising]")
			for m := range p.strategy.Out() {
				if p.l {
					trace.Log(p.name, "[utilised]", "v:", m.Value())
				}
				m.passP(p)
			}
		}
		return
	}
	for atomic.LoadInt64(&p.s)&stateClosed == 0x0 {
		select {
		case m, ok := <-p.strategy.Out():
			atomic.AddInt64(&p.s, -stateProcessing)
			if !ok {
				if p.l {
					trace.Log(p.name, "[sending broken value]")
				}
				atomic.AddInt64(&p.s, stateClosed)
				break
			}
			atomic.AddInt64(&p.s, stateReady)

			m.passP(p)
			if p.l {
				trace.Log(p.name, "[send]", "v:", m.Value())
			}
			p.out <- m

			atomic.AddInt64(&p.s, -stateReady)
			p.mu.Unlock()
		case <-p.ctx.Done():
			atomic.AddInt64(&p.s, stateClosed)
			if p.l {
				trace.Log(p.name, "[sending context deadline]")
			}
		}
	}
	close(p.out)
}

const (
	// stateProcessing means place is processing a token
	stateProcessing = 1 << 0
	// stateClosed means place is closed, and it doesn't process tokens
	stateClosed = 1 << 1
	// stateReady means place is ready to pass token forward
	stateReady = 1 << 2
)
