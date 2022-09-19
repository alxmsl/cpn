package cpn

import (
	"runtime"

	"github.com/alxmsl/cpn/trace"

	"github.com/alxmsl/prmtvs/skm"
)

// Transformation defines a custom behaviour for a transition
type Transformation func(in []*M) *M

// T implements an abstract transition in PN
type T struct {
	// name is a transition name in the PN. This is good to have it unique
	name string

	// transformation defines behaviour for the transition. Transition awaits tokens from each incoming edge. All tokens
	// are passed to the transformation. Transformation returns a token which will be passed to the following places
	transformation Transformation

	// ins is a sorted set of incoming edges
	ins *skm.SKM
	// outs is a sorted set of outgoing edges
	outs *skm.SKM

	// o keeps a static options flags for an abstract transition. See options constants for details
	o uint64
}

// NewT creates a new transition with required name
func NewT(name string) *T {
	var t = &T{
		name: name,

		ins:  skm.NewSKM(),
		outs: skm.NewSKM(),
	}
	if trace.NeedLog(t.name) {
		t.o &= optionLog
	}
	return t
}

func (t *T) SetOptions(opts ...TransitionOption) *T {
	for _, opt := range opts {
		opt.Apply(t)
	}
	return t
}

func (t *T) Name() string {
	return t.name
}

func (t *T) inslock() {
	t.ins.Over(func(i int, n string, v interface{}) bool {
		v.(*P).mu.Lock()
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
		v.(*P).mu.Unlock()
		return true
	})
}

func (t *T) run() {
	if t.o&optionLog > 0x0 {
		trace.Log(t.name, "[runinng...]")
		defer trace.Log(t.name, "[running completed]")
	}
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
		if t.o&optionLog > 0x0 {
			trace.Log(t.name, "[recv]", "len:", len(mm))
		}

		m := t.transformation(mm)
		m.passT(t)

		t.outs.Over(func(i int, n string, v interface{}) bool {
			in, _ := v.(*P).ins.GetByKey(t.Name())
			in.(chan *M) <- m
			return true
		})
		if t.o&optionLog > 0x0 {
			trace.Log(t.name, "[xfer]", "v:", m.Value())
		}
	}

	t.outs.Over(func(i int, n string, v interface{}) bool {
		in, _ := v.(*P).ins.GetByKey(t.Name())
		close(in.(chan *M))
		return true
	})
}
