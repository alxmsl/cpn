package cpn

import "github.com/alxmsl/prmtvs/skm"

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

func (pn *PN) PT(p, t string) *PN {
	pn.T(t).ins.Add(pn.P(p).Name(), pn.P(p))
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

func (pn *PN) TP(t, p string) *PN {
	pn.T(t).outs.Add(pn.P(p).Name(), pn.P(p))
	return pn
}

func (pn *PN) Run() {
	pn.pp.Over(func(i int, n string, v interface{}) bool {
		v.(*P).Run()
		return true
	})
	pn.tt.Over(func(i int, n string, v interface{}) bool {
		v.(*T).Run()
		return true
	})
}

func (pn *PN) Size() (int, int) {
	return pn.tt.Len(), pn.pp.Len()
}
