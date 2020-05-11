package test

import (
	"github.com/alxmsl/cpn/place"
	. "gopkg.in/check.v1"

	"context"
	"testing"

	"github.com/alxmsl/cpn"
	"github.com/alxmsl/cpn/transition"
)

func Test(t *testing.T) {
	TestingT(t)
}

type PNSuite struct{}

var _ = Suite(&PNSuite{})

func (s *PNSuite) TestPTP(c *C) {
	n := cpn.NewPN()
	n.P("pin", cpn.WithContext(context.Background()), cpn.WithPlace(place.NewBlock()))
	n.T("t1", cpn.WithFunction(transition.First))
	n.P("pout", cpn.WithContext(context.Background()), cpn.WithPlace(place.NewBlock()))

	n.PT("pin", "t1").TP("t1", "pout").Run()

	for i := 0; i < 1000; i += 1 {
		n.P("pin").In() <- cpn.NewM(i)
		m, ok := n.P("pout").Read()
		c.Assert(ok, Equals, true)
		c.Assert(m.Value().(int), Equals, i)
		c.Assert(m.Path(), HasLen, 3)
		c.Assert(m.Path()[0], Equals, "pin")
		c.Assert(m.Path()[1], Equals, "t1")
		c.Assert(m.Path()[2], Equals, "pout")
	}
}

func (s *PNSuite) TestPPTP(c *C) {
	n := cpn.NewPN()
	n.P("p1", cpn.WithContext(context.Background()), cpn.WithPlace(place.NewBlock()))
	n.P("p2", cpn.WithContext(context.Background()), cpn.WithPlace(place.NewBlock()))
	n.T("t1", cpn.WithFunction(transition.First))
	n.P("pout", cpn.WithContext(context.Background()), cpn.WithPlace(place.NewBlock()))

	n.PT("p1", "t1").PT("p2", "t1").TP("t1", "pout").Run()

	for i := 0; i < 1000; i += 1 {
		n.P("p1").In() <- cpn.NewM(i)
		n.P("p2").In() <- cpn.NewM(i)
		m, ok := n.P("pout").Read()
		c.Assert(ok, Equals, true)
		c.Assert(m.Value().(int), Equals, i)
		c.Assert(m.Path(), HasLen, 3)
		c.Assert(m.Path()[1], Equals, "t1")
		c.Assert(m.Path()[2], Equals, "pout")
	}
}

func (s *PNSuite) TestPPTTP(c *C) {
	n := cpn.NewPN()
	n.P("p1", cpn.WithContext(context.Background()), cpn.WithPlace(place.NewBlock()))
	n.P("p2", cpn.WithContext(context.Background()), cpn.WithPlace(place.NewBlock()))
	n.T("t1", cpn.WithFunction(transition.First))
	n.T("t2", cpn.WithFunction(transition.First))
	n.P("pout", cpn.WithContext(context.Background()), cpn.WithPlace(place.NewBlock()))

	n.PT("p1", "t1").PT("p2", "t1").PT("p1", "t2").PT("p2", "t2").
		TP("t1", "pout").TP("t2", "pout").Run()

	for i := 0; i < 1000; i += 1 {
		n.P("p1").In() <- cpn.NewM(i)
		n.P("p2").In() <- cpn.NewM(i)

		m, ok := n.P("pout").Read()
		c.Assert(ok, Equals, true)
		c.Assert(m.Value().(int), Equals, i)
	}
}
