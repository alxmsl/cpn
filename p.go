package rtpn

import (
	"context"
	"sync"
	"sync/atomic"
)

type Place interface {
	ReadCh() <-chan *M
	WriteCh() chan<- *M
}

type P struct {
	lock sync.Mutex

	ctx context.Context
	n   string

	in  chan *M
	out chan *M

	// Storage implementation
	place Place
	s     int64
}

func NewP(name string) *P {
	p := &P{
		n: name,

		in:  make(chan *M),
		out: make(chan *M),
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
	return p.n
}

func (p *P) Read() (*M, bool) {
	p.lock.Lock()
	m, ok := <-p.out
	return m, ok
}

func (p *P) Ready() bool {
	v := atomic.LoadInt64(&p.s)
	return v > 0
}

func (p *P) Run() {
	go func() {
		for m := range p.in {
			//p.storage <- m
			p.place.WriteCh() <- m
		}
	}()

	go func() {
		closed := false
		for !closed {
			select {
			//case m, ok := <-p.storage:
			case m, ok := <-p.place.ReadCh():
				if !ok {
					closed = true
					break
				}
				atomic.AddInt64(&p.s, 1)

				m.path = append(m.path, p.Name())
				p.out <- m

				atomic.AddInt64(&p.s, -1)
				p.lock.Unlock()
			case <-p.ctx.Done():
				closed = true
			}
		}
		close(p.out)
	}()
}

func (p *P) WriteCh() chan<- *M {
	return p.place.WriteCh()
}
