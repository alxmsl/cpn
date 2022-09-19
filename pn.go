package cpn

import (
	"fmt"
	"sync"

	"github.com/alxmsl/prmtvs/skm"
)

const formatName = "%s:%d"

type PN struct {
	pp *skm.SKM
	tt *skm.SKM
}

func NewPN() *PN {
	pn := &PN{
		pp: skm.NewSKM(),
		tt: skm.NewSKM(),
	}
	return pn
}

func (pn *PN) P(name string, opts ...PlaceOption) *P {
	if v, ok := pn.pp.GetByKey(name); ok {
		return v.(*P)
	}
	p := NewP(name).SetOptions(opts...)
	pn.pp.Add(p.Name(), p)
	return p
}

func (pn *PN) Pn(n int, prefix string, opts ...PlaceOption) {
	for i := 0; i < n; i += 1 {
		name := fmt.Sprintf(formatName, prefix, i)
		pn.P(name, opts...)
	}
}

func (pn *PN) PT(p, t string) *PN {
	pn.T(t).ins.Add(pn.P(p).Name(), pn.P(p))
	pn.P(p).o &= ^optionTerminal
	return pn
}

func (pn *PN) PTn(n int, p, prefix string) *PN {
	for i := 0; i < n; i += 1 {
		t := fmt.Sprintf(formatName, prefix, i)
		pn.PT(p, t)
	}
	return pn
}

func (pn *PN) PnTn(n int, prefixp, prefixt string) *PN {
	for i := 0; i < n; i += 1 {
		p := fmt.Sprintf(formatName, prefixp, i)
		t := fmt.Sprintf(formatName, prefixt, i)
		pn.PT(p, t)
	}
	return pn
}

func (pn *PN) PnT(n int, prefixp, t string) *PN {
	for i := 0; i < n; i += 1 {
		p := fmt.Sprintf(formatName, prefixp, i)
		pn.PT(p, t)
	}
	return pn
}

func (pn *PN) T(name string, opts ...TransitionOption) *T {
	if v, ok := pn.tt.GetByKey(name); ok {
		return v.(*T)
	}
	t := NewT(name).SetOptions(opts...)
	pn.tt.Add(t.Name(), t)
	return t
}

func (pn *PN) Tn(n int, prefix string, opts ...TransitionOption) {
	for i := 0; i < n; i += 1 {
		name := fmt.Sprintf(formatName, prefix, i)
		pn.T(name, opts...)
	}
}

func (pn *PN) TP(t, p string) *PN {
	pn.P(p).ins.Add(pn.T(t).Name(), make(chan *M))
	pn.T(t).outs.Add(pn.P(p).Name(), pn.P(p))
	pn.P(p).o &= ^optionInitial
	return pn
}

func (pn *PN) TPn(n int, t, prefixp string) *PN {
	for i := 0; i < n; i += 1 {
		p := fmt.Sprintf(formatName, prefixp, i)
		pn.TP(t, p)
	}
	return pn
}

func (pn *PN) TnPn(n int, prefixt, prefixp string) *PN {
	for i := 0; i < n; i += 1 {
		t := fmt.Sprintf(formatName, prefixt, i)
		p := fmt.Sprintf(formatName, prefixp, i)
		pn.TP(t, p)
	}
	return pn
}

func (pn *PN) TnP(n int, prefixt, p string) *PN {
	for i := 0; i < n; i += 1 {
		t := fmt.Sprintf(formatName, prefixt, i)
		pn.TP(t, p)
	}
	return pn
}

func (pn *PN) Run() {
	pn.pp.Over(func(i int, n string, v interface{}) bool {
		go v.(*P).run()
		go v.(*P).recv()
		go v.(*P).send()
		return true
	})
	pn.tt.Over(func(i int, n string, v interface{}) bool {
		go v.(*T).run()
		return true
	})
}

func (pn *PN) RunSync() {
	wg := sync.WaitGroup{}
	pn.pp.Over(func(i int, n string, v interface{}) bool {
		wg.Add(3)
		go func() {
			defer wg.Done()
			v.(*P).run()
		}()
		go func() {
			defer wg.Done()
			v.(*P).recv()
		}()
		go func() {
			defer wg.Done()
			v.(*P).send()
		}()
		return true
	})
	pn.tt.Over(func(i int, n string, v interface{}) bool {
		wg.Add(1)
		go func() {
			defer wg.Done()
			v.(*T).run()
		}()
		return true
	})
	wg.Wait()
}

func (pn *PN) Size() (int, int) {
	return pn.tt.Len(), pn.pp.Len()
}
