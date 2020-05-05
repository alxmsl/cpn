package test

import (
	. "gopkg.in/check.v1"

	"context"
	"testing"

	"github.com/alxmsl/rtpn"
	"github.com/alxmsl/rtpn/place"
	"github.com/alxmsl/rtpn/transition"
)

func Test(t *testing.T) {
	TestingT(t)
}

type PNSuite struct{}

var _ = Suite(&PNSuite{})

func (s *PNSuite) TestPTP(c *C) {
	pin := rtpn.NewP("pin").
		SetOptions(rtpn.WithContext(context.Background())).
		SetOptions(rtpn.WithPlace(place.NewBlock()))
	pout := rtpn.NewP("pout").
		SetOptions(rtpn.WithContext(context.Background())).
		SetOptions(rtpn.WithPlace(place.NewBlock()))
	t1 := rtpn.NewT("t1").
		SetOptions(rtpn.WithFunction(transition.First))

	n := rtpn.NewPN()
	n.PT(pin, t1)
	n.TP(t1, pout)
	n.Run()

	for i := 0; i < 1000; i += 1 {
		pin.WriteCh() <- rtpn.NewM(i)
		m, ok := pout.Read()
		c.Assert(ok, Equals, true)
		c.Assert(m.Value().(int), Equals, i)
		c.Assert(m.Path(), HasLen, 3)
		c.Assert(m.Path()[0], Equals, "pin")
		c.Assert(m.Path()[1], Equals, "t1")
		c.Assert(m.Path()[2], Equals, "pout")
	}
}

func (s *PNSuite) TestPPTP(c *C) {
	p1 := rtpn.NewP("p1").
		SetOptions(rtpn.WithContext(context.Background())).
		SetOptions(rtpn.WithPlace(place.NewBlock()))
	p2 := rtpn.NewP("p2").
		SetOptions(rtpn.WithContext(context.Background())).
		SetOptions(rtpn.WithPlace(place.NewBlock()))
	pout := rtpn.NewP("pout").
		SetOptions(rtpn.WithContext(context.Background())).
		SetOptions(rtpn.WithPlace(place.NewBlock()))
	t1 := rtpn.NewT("t1").
		SetOptions(rtpn.WithFunction(transition.First))

	n := rtpn.NewPN()
	n.PT(p1, t1)
	n.PT(p2, t1)
	n.TP(t1, pout)
	n.Run()

	for i := 0; i < 1000; i += 1 {
		p1.WriteCh() <- rtpn.NewM(i)
		p2.WriteCh() <- rtpn.NewM(i)
		m, ok := pout.Read()
		c.Assert(ok, Equals, true)
		c.Assert(m.Value().(int), Equals, i)
		c.Assert(m.Path(), HasLen, 3)
		c.Assert(m.Path()[1], Equals, "t1")
		c.Assert(m.Path()[2], Equals, "pout")
	}
}

func (s *PNSuite) TestPPTTP(c *C) {
	p1 := rtpn.NewP("p1").
		SetOptions(rtpn.WithContext(context.Background())).
		SetOptions(rtpn.WithPlace(place.NewBlock()))
	p2 := rtpn.NewP("p2").
		SetOptions(rtpn.WithContext(context.Background())).
		SetOptions(rtpn.WithPlace(place.NewBlock()))
	pout := rtpn.NewP("pout").
		SetOptions(rtpn.WithContext(context.Background())).
		SetOptions(rtpn.WithPlace(place.NewBlock()))
	t1 := rtpn.NewT("t1").
		SetOptions(rtpn.WithFunction(transition.First))
	t2 := rtpn.NewT("t2").
		SetOptions(rtpn.WithFunction(transition.First))

	n := rtpn.NewPN()
	n.PT(p1, t1)
	n.PT(p2, t1)
	n.PT(p1, t2)
	n.PT(p2, t2)
	n.TP(t1, pout)
	n.TP(t2, pout)
	n.Run()

	for i := 0; i < 1000; i += 1 {
		p1.WriteCh() <- rtpn.NewM(i)
		p2.WriteCh() <- rtpn.NewM(i)

		m, ok := pout.Read()
		c.Assert(ok, Equals, true)
		c.Assert(m.Value().(int), Equals, i)
	}
}
