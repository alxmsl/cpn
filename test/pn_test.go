package test

import (
	. "gopkg.in/check.v1"

	"context"
	"testing"

	"github.com/alxmsl/pn"
	"github.com/alxmsl/pn/place"
)

func Test(t *testing.T) {
	TestingT(t)
}

type PNSuite struct{}

var _ = Suite(&PNSuite{})

func (s *PNSuite) TestPTP(c *C) {
	pin := pn.NewP("pin").
		SetOptions(pn.WithContextOption(context.Background())).
		SetOptions(pn.WithPlaceOption(place.NewBlock()))
	pout := pn.NewP("pout").
		SetOptions(pn.WithContextOption(context.Background())).
		SetOptions(pn.WithPlaceOption(place.NewBlock()))
	t1 := pn.NewT("t1")

	n := pn.NewPN()
	n.PT(pin, t1)
	n.TP(t1, pout)
	n.Run()

	for i := 0; i < 1000; i += 1 {
		pin.WriteCh() <- pn.NewM(i)
		m, ok := pout.Read()
		c.Assert(ok, Equals, true)
		c.Assert(m.Value(), Equals, i)
		c.Assert(m.Path(), HasLen, 3)
		c.Assert(m.Path()[0], Equals, "pin")
		c.Assert(m.Path()[1], Equals, "t1")
		c.Assert(m.Path()[2], Equals, "pout")
	}
}

func (s *PNSuite) TestPPTP(c *C) {
	p1 := pn.NewP("p1").
		SetOptions(pn.WithContextOption(context.Background())).
		SetOptions(pn.WithPlaceOption(place.NewBlock()))
	p2 := pn.NewP("p2").
		SetOptions(pn.WithContextOption(context.Background())).
		SetOptions(pn.WithPlaceOption(place.NewBlock()))
	pout := pn.NewP("pout").
		SetOptions(pn.WithContextOption(context.Background())).
		SetOptions(pn.WithPlaceOption(place.NewBlock()))
	t1 := pn.NewT("t1")

	n := pn.NewPN()
	n.PT(p1, t1)
	n.PT(p2, t1)
	n.TP(t1, pout)
	n.Run()

	for i := 0; i < 1000; i += 1 {
		p1.WriteCh() <- pn.NewM(i)
		p2.WriteCh() <- pn.NewM(i)
		m, ok := pout.Read()
		c.Assert(ok, Equals, true)
		c.Assert(m.Value(), Equals, i*2)
		c.Assert(m.Path(), HasLen, 3)
		c.Assert(m.Path()[1], Equals, "t1")
		c.Assert(m.Path()[2], Equals, "pout")
	}
}

func (s *PNSuite) TestPPTTP(c *C) {
	p1 := pn.NewP("p1").
		SetOptions(pn.WithContextOption(context.Background())).
		SetOptions(pn.WithPlaceOption(place.NewBlock()))
	p2 := pn.NewP("p2").
		SetOptions(pn.WithContextOption(context.Background())).
		SetOptions(pn.WithPlaceOption(place.NewBlock()))
	pout := pn.NewP("pout").
		SetOptions(pn.WithContextOption(context.Background())).
		SetOptions(pn.WithPlaceOption(place.NewBlock()))
	t1 := pn.NewT("t1")
	t2 := pn.NewT("t2")

	n := pn.NewPN()
	n.PT(p1, t1)
	n.PT(p2, t1)
	n.PT(p1, t2)
	n.PT(p2, t2)
	n.TP(t1, pout)
	n.TP(t2, pout)
	n.Run()

	for i := 0; i < 1000; i += 1 {
		p1.WriteCh() <- pn.NewM(i)
		p2.WriteCh() <- pn.NewM(i)

		m, ok := pout.Read()
		c.Assert(ok, Equals, true)
		c.Assert(m.Value(), Equals, i*2)
	}
}
