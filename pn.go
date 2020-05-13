package cpn

import (
	"fmt"
	"sync"

	"github.com/alxmsl/prmtvs/skm"
)

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

func (pn *PN) P(name string, opts ...POpt) *P {
	if v, ok := pn.pp.GetByKey(name); ok {
		return v.(*P)
	}
	p := NewP(name).SetOptions(opts...)
	pn.pp.Add(p.Name(), p)
	return p
}

func (pn *PN) Pn(n int, prefix string, opts ...POpt) {
	for i := 0; i < n; i += 1 {
		name := fmt.Sprintf("%s%d", prefix, i)
		pn.P(name, opts...)
	}
}

func (pn *PN) PT(p, t string) *PN {
	pn.T(t).ins.Add(pn.P(p).Name(), pn.P(p))
	return pn
}

func (pn *PN) PTn(n int, p, prefix string) *PN {
	for i := 0; i < n; i += 1 {
		t := fmt.Sprintf("%s%d", prefix, i)
		pn.PT(p, t)
	}
	return pn
}

func (pn *PN) T(name string, opts ...TOpt) *T {
	if v, ok := pn.tt.GetByKey(name); ok {
		return v.(*T)
	}
	t := NewT(name).SetOptions(opts...)
	pn.tt.Add(t.Name(), t)
	return t
}

func (pn *PN) Tn(n int, prefix string, opts ...TOpt) {
	for i := 0; i < n; i += 1 {
		name := fmt.Sprintf("%s%d", prefix, i)
		pn.T(name, opts...)
	}
}

func (pn *PN) TP(t, p string) *PN {
	pn.T(t).outs.Add(pn.P(p).Name(), pn.P(p))
	return pn
}

func (pn *PN) TnPn(n int, prefixt, prefixp string) *PN {
	for i := 0; i < n; i += 1 {
		t := fmt.Sprintf("%s%d", prefixt, i)
		p := fmt.Sprintf("%s%d", prefixp, i)
		pn.TP(t, p)
	}
	return pn
}

func (pn *PN) Run() {
	pn.pp.Over(func(i int, n string, v interface{}) bool {
		go v.(*P).startPlace()
		go v.(*P).startRecv()
		go v.(*P).startSend()
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
			v.(*P).startPlace()
		}()
		go func() {
			defer wg.Done()
			v.(*P).startRecv()
		}()
		go func() {
			defer wg.Done()
			v.(*P).startSend()
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
