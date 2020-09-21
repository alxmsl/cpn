package cpn

import (
	"runtime"

	"github.com/alxmsl/prmtvs/skm"
)

type Transition func(in []*M) *M

type T struct {
	n string

	ins  *skm.SKM
	outs *skm.SKM

	fn Transition
}

func NewT(name string) *T {
	t := &T{
		n: name,

		ins:  skm.NewSKM(),
		outs: skm.NewSKM(),
	}
	return t
}

func (t *T) SetOptions(opts ...TOpt) *T {
	for _, opt := range opts {
		opt.Apply(t)
	}
	return t
}

func (t *T) Name() string {
	return t.n
}

func (t *T) inslock() {
	t.ins.Over(func(i int, n string, v interface{}) bool {
		v.(*P).lock.Lock()
		return true
	})
}

func (t *T) insready() bool {
	var ready bool
	t.ins.Over(func(i int, n string, v interface{}) bool {
		if i > 0 {
			ready = ready && v.(*P).ready()
		} else {
			ready = v.(*P).ready()
		}
		return ready
	})
	return ready
}

func (t *T) insunlock() {
	t.ins.Over(func(i int, n string, v interface{}) bool {
		v.(*P).lock.Unlock()
		return true
	})
}

func (t *T) run() {
	for {
		t.inslock()
		if !t.insready() {
			t.insunlock()
			runtime.Gosched()
			continue
		}

		mm := make([]*M, t.ins.Len())
		var ok bool
		t.ins.Over(func(i int, n string, v interface{}) bool {
			mm[i], ok = <-v.(*P).out
			return ok
		})
		if !ok {
			t.insunlock()
			break
		}

		m := t.fn(mm)
		if m != nil {
			continue
		}
		m.PassT(t.Name())

		t.outs.Over(func(i int, n string, v interface{}) bool {
			in, _ := v.(*P).ins.GetByKey(t.Name())
			in.(chan *M) <- m
			return true
		})
	}

	t.outs.Over(func(i int, n string, v interface{}) bool {
		in, _ := v.(*P).ins.GetByKey(t.Name())
		close(in.(chan *M))
		return true
	})
}
