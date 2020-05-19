package cpn

import (
	"context"
	"sync"
	"sync/atomic"
)

type Place interface {
	In() chan<- *M
	Out() <-chan *M
	Close()
	Run()
}

const (
	stateActive = 1 << 0
	stateClosed = 1 << 1
	stateReady  = 1 << 2
)

type P struct {
	lock sync.Mutex

	ctx context.Context
	n   string

	in  chan *M
	out chan *M

	// Storage implementation
	place Place
	s     int64
	i     bool
	t     bool
}

func NewP(name string) *P {
	p := &P{
		n: name,

		in:  make(chan *M),
		out: make(chan *M),
	}
	return p
}

func (p *P) SetOptions(opts ...POpt) *P {
	for _, opt := range opts {
		opt.Apply(p)
	}
	return p
}

func (p *P) Name() string {
	return p.n
}

func (p *P) In() chan<- *M {
	return p.place.In()
}

func (p *P) Out() <-chan *M {
	if !p.t {
		return nil
	}
	return p.place.Out()
}

func (p *P) ready() bool {
	s := atomic.LoadInt64(&p.s)
	return s&(stateReady|stateClosed) > 0x0
}

func (p *P) startPlace() {
	p.place.Run()
}

func (p *P) startRecv() {
	atomic.AddInt64(&p.s, stateActive)
	if p.i {
		return
	}
	for m := range p.in {
		m.path = append(m.path, p.Name())
		p.place.In() <- m
	}
	if p.t {
		p.place.Close()
		close(p.out)
	}
}

func (p *P) startSend() {
	if p.t {
		return
	}
	for atomic.LoadInt64(&p.s)&stateClosed == 0x0 {
		select {
		case m, ok := <-p.place.Out():
			if !ok {
				atomic.AddInt64(&p.s, stateClosed)
				break
			}
			atomic.AddInt64(&p.s, stateReady)

			if len(m.path) == 0 || m.path[len(m.path)-1] != p.Name() {
				m.path = append(m.path, p.Name())
			}
			p.out <- m

			atomic.AddInt64(&p.s, -stateReady)
			p.lock.Unlock()
		case <-p.ctx.Done():
			atomic.AddInt64(&p.s, stateClosed)
		}
	}
	close(p.out)
}
