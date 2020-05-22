package test

import (
	. "gopkg.in/check.v1"

	"bytes"
	"context"
	"testing"

	"github.com/alxmsl/cpn"
	"github.com/alxmsl/cpn/place"
	"github.com/alxmsl/cpn/transition"
)

func Test(t *testing.T) {
	TestingT(t)
}

type PNSuite struct{}

var _ = Suite(&PNSuite{})

func (s *PNSuite) TestPTP(c *C) {
	ctx, cancel := context.WithCancel(context.Background())

	n := cpn.NewPN()
	n.P("pin", cpn.WithContext(ctx), cpn.WithPlace(place.NewBlock()), cpn.IsInitial())
	n.T("t1", cpn.WithFunction(transition.First))
	n.P("pout", cpn.WithContext(ctx), cpn.WithPlace(place.NewBlock()), cpn.IsFinal())

	n.PT("pin", "t1").TP("t1", "pout").Run()

	go func() {
		for i := 0; i < 1000; i += 1 {
			n.P("pin").In() <- cpn.NewM(i)
		}
		cancel()
	}()

	i := 0
	for m := range n.P("pout").Out() {
		c.Assert(m.Value().(int), Equals, i)
		c.Assert(m.Path(), HasLen, 3)
		c.Assert(m.Path()[0], Equals, "pin")
		c.Assert(m.Path()[1], Equals, "t1")
		c.Assert(m.Path()[2], Equals, "pout")
		i += 1
	}
}

func (s *PNSuite) TestPTTP(c *C) {
	ctx, cancel := context.WithCancel(context.Background())

	n := cpn.NewPN()
	n.P("pin", cpn.WithContext(ctx), cpn.WithPlace(place.NewBlock()), cpn.IsInitial())
	n.T("t1", cpn.WithFunction(transition.First))
	n.T("t2", cpn.WithFunction(transition.First))
	n.P("pout", cpn.WithContext(ctx), cpn.WithPlace(place.NewBlock()), cpn.IsFinal())

	n.PT("pin", "t1").PT("pin", "t2").TP("t1", "pout").TP("t2", "pout").Run()

	go func() {
		for i := 0; i < 1000; i += 1 {
			n.P("pin").In() <- cpn.NewM(i)
		}
		cancel()
	}()

	count := 0
	for m := range n.P("pout").Out() {
		c.Assert(m.Path(), HasLen, 3)
		c.Assert(m.Path()[0], Equals, "pin")
		c.Assert(m.Path()[1] == "t1" || m.Path()[1] == "t2", Equals, true)
		c.Assert(m.Path()[2], Equals, "pout")
		count += 1
	}
	c.Assert(count, Equals, 1000)
}

func (s *PNSuite) TestPPTP(c *C) {
	ctx, cancel := context.WithCancel(context.Background())

	n := cpn.NewPN()
	n.P("p1", cpn.WithContext(ctx), cpn.WithPlace(place.NewBlock()), cpn.IsInitial())
	n.P("p2", cpn.WithContext(ctx), cpn.WithPlace(place.NewBlock()), cpn.IsInitial())
	n.T("t1", cpn.WithFunction(transition.First))
	n.P("pout", cpn.WithContext(ctx), cpn.WithPlace(place.NewBlock()), cpn.IsFinal())

	n.PT("p1", "t1").PT("p2", "t1").TP("t1", "pout").Run()

	go func() {
		for i := 0; i < 1000; i += 1 {
			n.P("p1").In() <- cpn.NewM(i)
			n.P("p2").In() <- cpn.NewM(i)
		}
		cancel()
	}()

	i := 0
	for m := range n.P("pout").Out() {
		c.Assert(m.Value().(int), Equals, i)
		c.Assert(m.Path(), HasLen, 3)
		c.Assert(m.Path()[0] == "p1" || m.Path()[0] == "p2", Equals, true)
		c.Assert(m.Path()[1], Equals, "t1")
		c.Assert(m.Path()[2], Equals, "pout")
		i += 1
	}
}

func (s *PNSuite) TestPPTTP(c *C) {
	ctx, cancel := context.WithCancel(context.Background())

	n := cpn.NewPN()
	n.P("p1", cpn.WithContext(ctx), cpn.WithPlace(place.NewBlock()), cpn.IsInitial())
	n.P("p2", cpn.WithContext(ctx), cpn.WithPlace(place.NewBlock()), cpn.IsInitial())
	n.T("t1", cpn.WithFunction(transition.First))
	n.T("t2", cpn.WithFunction(transition.First))
	n.P("pout", cpn.WithContext(ctx), cpn.WithPlace(place.NewBlock()), cpn.IsFinal())

	n.PT("p1", "t1").PT("p2", "t1").PT("p1", "t2").PT("p2", "t2").
		TP("t1", "pout").TP("t2", "pout").Run()

	go func() {
		for i := 0; i < 1000; i += 1 {
			n.P("p1").In() <- cpn.NewM(i)
			n.P("p2").In() <- cpn.NewM(i)
		}
		cancel()
	}()

	count := 0
	for m := range n.P("pout").Out() {
		c.Assert(m.Path(), HasLen, 3)
		c.Assert(m.Path()[0] == "p1" || m.Path()[0] == "p2", Equals, true)
		c.Assert(m.Path()[1] == "t1" || m.Path()[1] == "t2", Equals, true)
		c.Assert(m.Path()[2], Equals, "pout")
		count += 1
	}
	c.Assert(count, Equals, 1000)
}

func (s *PNSuite) TestPrintNet(c *C) {
	ctx, cancel := context.WithCancel(context.Background())
	n := cpn.NewPN()
	n.P("pin", cpn.WithContext(ctx), cpn.WithPlace(place.NewBlock()), cpn.IsInitial())
	n.T("t1", cpn.WithFunction(transition.First))

	w := bytes.NewBufferString("")
	n.P("pout", cpn.WithContext(ctx), cpn.WithPlace(place.NewPrint(place.WriterOption(w))), cpn.IsFinal())

	n.PT("pin", "t1").TP("t1", "pout")

	go func() {
		for i := 0; i < 10; i += 1 {
			n.P("pin").In() <- cpn.NewM(i)
		}
		cancel()
	}()
	n.RunSync()
	c.Assert(w.String(), Equals, "0\n1\n2\n3\n4\n5\n6\n7\n8\n9\n")
}
