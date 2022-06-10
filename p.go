package cpn

import (
	"context"
	"log"
	"sync"
	"sync/atomic"

	"github.com/alxmsl/prmtvs/skm"
)

type Place interface {
	In() chan<- *M
	Out() <-chan *M
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

	ch  chan struct{}
	in  chan *M
	ins *skm.SKM
	out chan *M

	// Storage implementation
	place   Place
	s       int64
	i, k, t bool
}

func NewP(name string) *P {
	p := &P{
		n: name,

		ch:  make(chan struct{}),
		in:  make(chan *M),
		ins: skm.NewSKM(),
		out: make(chan *M),

		i: true,
		t: true,
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

func (p *P) Close() {
	close(p.place.In())
}

func (p *P) In() chan<- *M {
	return p.in
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
	log.Println("startPlace", p.n)
	p.place.Run()
}

func (p *P) startRecv() {
	log.Println("startRecv", p.n)
	atomic.AddInt64(&p.s, stateActive)
	if p.i {
		go func() {
			for m := range p.in {
				m.PassP(p.Name())
				p.place.In() <- m
			}
			close(p.ch)
		}()
		return
	}

	wg := &sync.WaitGroup{}
	wg.Add(p.ins.Len())
	p.ins.Over(func(i int, n string, v interface{}) bool {
		go func() {
			defer wg.Done()
			for m := range v.(chan *M) {
				m.PassP(p.Name())
				p.place.In() <- m
			}
		}()
		return true
	})
	wg.Wait()

	close(p.place.In())
	if p.t {
		close(p.out)
	}
}

func (p *P) startSend() {
	log.Println("startsend", p.n)
	if p.t {
		if !p.k {
			for range p.place.Out() {
			}
		}
		return
	}
	for atomic.LoadInt64(&p.s)&stateClosed == 0x0 {
		select {
		case m, ok := <-p.place.Out():
			log.Println("2", p.n)
			if !ok {
				atomic.AddInt64(&p.s, stateClosed)
				break
			}
			atomic.AddInt64(&p.s, stateReady)

			m.PassP(p.Name())
			p.out <- m

			atomic.AddInt64(&p.s, -stateReady)
			log.Println("unlock", p.n)
			p.lock.Unlock()
		case <-p.ctx.Done():
			log.Println("3", p.n)
			close(p.in)
			if p.i {
				<-p.ch
			}
			atomic.AddInt64(&p.s, stateClosed)
		}
	}
	log.Println("close", p.n)
	close(p.out)
}
