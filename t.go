package pn

import (
	"runtime"

	"github.com/alxmsl/prmtvs/skm"
)

type T struct {
	n string

	ins  *skm.SKM
	outs *skm.SKM
}

func NewT(name string) *T {
	t := &T{
		n: name,

		ins:  skm.NewSKM(),
		outs: skm.NewSKM(),
	}
	return t
}

func (t *T) Name() string {
	return t.n
}

func (t *T) Run() {
	go func() {
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
				break
			}

			m := mm[0]
			for i := 1; i < t.ins.Len(); i += 1 {
				m.v += mm[i].v
			}
			m.path = append(m.path, t.Name())
			m.word = append(m.word, t.Name())

			t.outs.Over(func(i int, n string, v interface{}) bool {
				v.(*P).in <- m
				return true
			})

		}

		t.outs.Over(func(i int, n string, v interface{}) bool {
			close(v.(*P).in)
			return true
		})
	}()
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
			ready = ready && v.(*P).Ready()
		} else {
			ready = v.(*P).Ready()
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
