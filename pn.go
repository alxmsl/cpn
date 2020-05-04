package rtpn

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

func (pn *PN) PT(p *P, t *T) *PN {
	pn.pp.Add(p.Name(), p)
	pn.tt.Add(t.Name(), t)

	t.ins.Add(p.Name(), p)
	return pn
}

func (pn *PN) TP(t *T, p *P) *PN {
	pn.pp.Add(p.Name(), p)
	pn.tt.Add(t.Name(), t)

	t.outs.Add(p.Name(), p)
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
