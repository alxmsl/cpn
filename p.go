package cpn

import (
	"context"
	"sync"

	"github.com/alxmsl/cpn/trace"
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
	ctx  context.Context
	lock sync.RWMutex
	mu   sync.Mutex

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

	// o keeps a static options flags for an abstract place. See options constants for details
	o uint64

	// s keeps a dynamic state for the place. See state constants for details
	s uint64
}

// NewP creates a new place with specific name. By default, place is both initial and terminal. When place is linked to
// transition default state is changed
func NewP(name string) *P {
	var p = &P{
		name: name,

		ins: skm.NewSKM(),
		out: make(chan *M),

		o: optionInitial | optionTerminal,
	}
	if trace.NeedLog(p.name) {
		p.o |= optionLog
	}
	return p
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
	if p.o&optionTerminal > 0x0 {
		return p.strategy.Out()
	}
	return nil
}

func (p *P) Send(m *M) {
	if p.o&optionLog > 0x0 {
		trace.Log(p.name, "[recv (direct)]", "v:", m.Value())
	}
	p.setState(stateProcessing)
	p.In() <- m
}

func (p *P) setState(state uint64) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.s |= state
}

func (p *P) unsetState(state uint64) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.s &= ^state
}

func (p *P) state() uint64 {
	p.lock.RLock()
	defer p.lock.RUnlock()
	return p.s
}

func (p *P) ready() bool {
	return p.state()&(stateReady|stateClosed) > 0x0
}

func (p *P) run() {
	if p.o&optionLog > 0x0 {
		trace.Logf("%s [running...] o:%064b\n", p.name, p.o)
		defer trace.Log(p.name, "[running completed]")
	}
	p.strategy.Run()
}

func (p *P) recv() {
	if p.o&optionLog > 0x0 {
		trace.Log(p.name, "[receiving...]")
		defer trace.Log(p.name, "[receiving completed]")
	}

	if p.o&optionInitial > 0x0 {
		if p.o&optionLog > 0x0 {
			trace.Log(p.name, "[initial place detected]")
		}
		return
	}

	wg := &sync.WaitGroup{}
	wg.Add(p.ins.Len())
	p.ins.Over(func(i int, n string, v interface{}) bool {
		if p.o&optionLog > 0x0 {
			trace.Log(p.name, "[listening...]", "n:", n)
		}
		go func() {
			if p.o&optionLog > 0x0 {
				defer trace.Log(p.name, "[stop listening]", "n:", n)
			}
			defer wg.Done()
			for m := range v.(chan *M) {
				m.passP(p)
				if p.o&optionLog > 0x0 {
					trace.Log(p.name, "[recv]", "n:", n, "v:", m.Value())
				}
				p.setState(stateProcessing)
				p.In() <- m
			}
		}()
		return true
	})
	wg.Wait()

	close(p.strategy.In())
	if p.o&optionTerminal > 0x0 {
		close(p.out)
	}
}

func (p *P) send() {
	if p.o&optionLog > 0x0 {
		trace.Log(p.name, "[sending...]")
		defer trace.Log(p.name, "[sending completed]")
	}
	if p.o&optionTerminal > 0x0 {
		if p.o&optionLog > 0x0 {
			trace.Log(p.name, "[terminal place detected]")
		}
		if p.o&optionKeep > 0x0 {
			return
		}
		if p.o&optionLog > 0x0 {
			trace.Log(p.name, "[utilising]")
		}
		for m := range p.strategy.Out() {
			if p.o&optionLog > 0x0 {
				trace.Log(p.name, "[utilised]", "v:", m.Value())
			}
			m.passP(p)
		}
		return
	}
	for s := p.state(); (s&stateClosed)|(^s&stateProcessing) != (stateClosed | stateProcessing); s = p.state() {
		select {
		case m, ok := <-p.strategy.Out():
			if !ok {
				if p.o&optionLog > 0x0 {
					trace.Log(p.name, "[sending broken value]")
				}
				p.setState(stateClosed)
				p.unsetState(stateProcessing)
				break
			}
			p.setState(stateReady)
			p.unsetState(stateProcessing)

			m.passP(p)
			if p.o&optionLog > 0x0 {
				trace.Log(p.name, "[send]", "v:", m.Value())
			}
			p.out <- m

			p.unsetState(stateReady)
			p.mu.Unlock()
		case <-p.ctx.Done():
			p.setState(stateClosed)
			if p.o&optionLog > 0x0 {
				trace.Log(p.name, "[sending context deadline]")
			}
		}
	}
	close(p.out)
}

const (
	// stateClosed means place is closed, and it doesn't process tokens
	stateClosed uint64 = 1 << 0
	// stateProcessing means place is processing a token
	stateProcessing uint64 = 1 << 1
	// stateReady means place is ready to pass token forward
	stateReady uint64 = 1 << 2
)
