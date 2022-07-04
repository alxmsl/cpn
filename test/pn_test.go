package test

import (
	. "gopkg.in/check.v1"

	"bytes"
	"context"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/alxmsl/cpn"
	"github.com/alxmsl/cpn/place"
	"github.com/alxmsl/cpn/place/io"
	"github.com/alxmsl/cpn/place/memory"
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
	n.P("pin",
		cpn.WithContext(ctx),
		cpn.WithPlace(memory.NewBlock()),
	)
	n.T("t1", cpn.WithFunction(transition.First))
	n.P("pout",
		cpn.WithContext(ctx),
		cpn.WithPlace(memory.NewBlock()),
		cpn.WithKeep(true),
	)
	n.
		PT("pin", "t1").
		TP("t1", "pout").
		Run()
	go func() {
		for i := 0; i < 1000; i += 1 {
			n.P("pin").In() <- cpn.NewM(i)
		}
		cancel()
	}()

	i := 0
	for m := range n.P("pout").Out() {
		c.Assert(m.Value(), NotNil, Commentf("%#v", m))
		c.Assert(m.Value().(int), Equals, i)
		c.Assert(m.Path(), HasLen, 3)
		c.Assert(m.Path()[0].N, Equals, "pin")
		c.Assert(m.Path()[1].N, Equals, "t1")
		c.Assert(m.Path()[2].N, Equals, "pout")

		c.Assert(m.History(), HasLen, 4)
		c.Assert(m.History()[0].N, Equals, "")
		c.Assert(m.History()[1], Equals, m.Path()[0])
		c.Assert(m.History()[2], Equals, m.Path()[1])
		c.Assert(m.History()[3], Equals, m.Path()[2])

		c.Assert(m.ValueByPlace("pin", 0), NotNil)
		// Check that marks have been processed sequentially in the Petri Net
		c.Assert(m.ValueByPlace("pin", 0).(int), Equals, i)
		c.Assert(m.ValueByPlace("pout", 0), IsNil)

		i += 1
	}
}

func (s *PNSuite) TestPTTP(c *C) {
	ctx, cancel := context.WithCancel(context.Background())

	n := cpn.NewPN()
	n.P("pin",
		cpn.WithContext(ctx),
		cpn.WithPlace(memory.NewBlock()),
	)
	n.T("t1", cpn.WithFunction(transition.First))
	n.T("t2", cpn.WithFunction(transition.First))
	n.P("pout",
		cpn.WithContext(ctx),
		cpn.WithPlace(memory.NewBlock()),
		cpn.WithKeep(true),
	)
	n.
		PT("pin", "t1").
		PT("pin", "t2").
		TP("t1", "pout").
		TP("t2", "pout").
		Run()
	go func() {
		for i := 0; i < 1000; i += 1 {
			n.P("pin").In() <- cpn.NewM(i)
		}
		cancel()
	}()

	count := 0
	for m := range n.P("pout").Out() {
		c.Assert(m.Path(), HasLen, 3)
		c.Assert(m.Path()[0].N, Equals, "pin")
		c.Assert(m.Path()[1].N == "t1" || m.Path()[1].N == "t2", Equals, true)
		c.Assert(m.Path()[2].N, Equals, "pout")

		c.Assert(m.History(), HasLen, 4)
		c.Assert(m.History()[0].N, Equals, "")
		c.Assert(m.History()[1], Equals, m.Path()[0])
		c.Assert(m.History()[2], Equals, m.Path()[1])
		c.Assert(m.History()[3], Equals, m.Path()[2])

		// We cannot check the mark's value in sequence, because marks could be processed in different sequences: using
		// `t1` or `t2`. So, we check just that mark has value from place `pin` and doesn't have it from place `pout`
		c.Assert(m.ValueByPlace("pin", 0), NotNil)
		c.Assert(m.ValueByPlace("pout", 0), IsNil)

		count += 1
	}
	c.Assert(count, Equals, 1000)
}

func (s *PNSuite) TestPTPP(c *C) {
	ctx, cancel := context.WithCancel(context.Background())

	n := cpn.NewPN()
	n.P("pin",
		cpn.WithContext(ctx),
		cpn.WithPlace(memory.NewBlock()),
	)
	n.T("t", cpn.WithFunction(transition.First))
	n.P("pout1",
		cpn.WithContext(ctx),
		cpn.WithPlace(memory.NewBlock()),
		cpn.WithKeep(true),
	)
	n.P("pout2",
		cpn.WithContext(ctx),
		cpn.WithPlace(memory.NewBlock()),
		cpn.WithKeep(true),
	)
	n.
		PT("pin", "t").
		TP("t", "pout1").
		TP("t", "pout2").
		Run()
	go func() {
		for i := 0; i < 1000; i += 1 {
			n.P("pin").In() <- cpn.NewM(i)
		}
		cancel()
	}()

	var count int64
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		for m := range n.P("pout1").Out() {
			c.Assert(m.ValueByPlace("pin", 0), NotNil)
			c.Assert(m.ValueByPlace("pout1", 0), IsNil)
			c.Assert(m.ValueByPlace("pout2", 0), IsNil)

			atomic.AddInt64(&count, 1)
		}
	}()
	go func() {
		defer wg.Done()
		for m := range n.P("pout2").Out() {
			c.Assert(m.ValueByPlace("pin", 0), NotNil)
			c.Assert(m.ValueByPlace("pout1", 0), IsNil)
			c.Assert(m.ValueByPlace("pout2", 0), IsNil)

			atomic.AddInt64(&count, 1)
		}
	}()
	wg.Wait()
	c.Assert(count, Equals, int64(2000))
}

func (s *PNSuite) TestPPTP(c *C) {
	ctx, cancel := context.WithCancel(context.Background())

	n := cpn.NewPN()
	n.P("p1",
		cpn.WithContext(ctx),
		cpn.WithPlace(memory.NewBlock()),
	)
	n.P("p2",
		cpn.WithContext(ctx),
		cpn.WithPlace(memory.NewBlock()),
	)
	n.T("t1", cpn.WithFunction(transition.First))
	n.P("pout",
		cpn.WithContext(ctx),
		cpn.WithPlace(memory.NewBlock()),
		cpn.WithKeep(true),
	)
	n.
		PT("p1", "t1").
		PT("p2", "t1").
		TP("t1", "pout").
		Run()
	go func() {
		for i := 0; i < 1000; i += 1 {
			n.P("p1").In() <- cpn.NewM(i)
			n.P("p2").In() <- cpn.NewM(i)
		}
		cancel()
	}()

	i := 0
	for m := range n.P("pout").Out() {
		c.Assert(m.Value(), NotNil)
		c.Assert(m.Value().(int), Equals, i)
		c.Assert(m.Path(), HasLen, 3)
		// Path always has a place `p1`, because transition `t1` passes mark from the first place only
		c.Assert(m.Path()[0].N, Equals, "p1")
		c.Assert(m.Path()[1].N, Equals, "t1")
		c.Assert(m.Path()[2].N, Equals, "pout")

		c.Assert(m.History(), HasLen, 4)
		c.Assert(m.History()[0].N, Equals, "")
		c.Assert(m.History()[1], Equals, m.Path()[0])
		c.Assert(m.History()[2], Equals, m.Path()[1])
		c.Assert(m.History()[3], Equals, m.Path()[2])

		// Mark has value from place `p1` and doesn't have from place `p2` because transition `t1` passes mark from the
		// first place only
		c.Assert(m.ValueByPlace("p1", 0), NotNil)
		c.Assert(m.ValueByPlace("p1", 0).(int), Equals, i)
		c.Assert(m.ValueByPlace("p2", 0), IsNil)
		c.Assert(m.ValueByPlace("pout", 0), IsNil)

		i += 1
	}
}

func (s *PNSuite) TestPPTTP(c *C) {
	ctx, cancel := context.WithCancel(context.Background())

	n := cpn.NewPN()
	n.P("p1",
		cpn.WithContext(ctx),
		cpn.WithPlace(memory.NewBlock()),
	)
	n.P("p2",
		cpn.WithContext(ctx),
		cpn.WithPlace(memory.NewBlock()),
	)
	n.T("t1", cpn.WithFunction(transition.First))
	n.T("t2", cpn.WithFunction(transition.First))
	n.P("pout",
		cpn.WithContext(ctx),
		cpn.WithPlace(memory.NewBlock()),
		cpn.WithKeep(true),
	)
	n.
		PT("p1", "t1").
		PT("p2", "t1").
		PT("p1", "t2").
		PT("p2", "t2").
		TP("t1", "pout").
		TP("t2", "pout").
		Run()
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
		// Path always has a place `p1`, because both transitions `t1` and `t2` passes mark from the first place only.
		c.Assert(m.Path()[0].N, Equals, "p1")
		// Path can have both transitions here, because of concurrency execution.
		c.Assert(m.Path()[1].N == "t1" || m.Path()[1].N == "t2", Equals, true)
		c.Assert(m.Path()[2].N, Equals, "pout")

		c.Assert(m.History(), HasLen, 4)
		c.Assert(m.History()[0].N, Equals, "")
		c.Assert(m.History()[1], Equals, m.Path()[0])
		c.Assert(m.History()[2], Equals, m.Path()[1])
		c.Assert(m.History()[3], Equals, m.Path()[2])

		// Mark has value from place `p1` and doesn't have from place `p2` because both transitions `t1` and `t2` passes
		// mark from the first place only. The same moment we can not predict the mark's value, because both transitions
		// `t1` and `t2` work concurrently
		c.Assert(m.ValueByPlace("p1", 0), NotNil)
		c.Assert(m.ValueByPlace("p2", 0), IsNil)
		c.Assert(m.ValueByPlace("pout", 0), IsNil)

		count += 1
	}
	c.Assert(count, Equals, 1000)
}

func (s *PNSuite) TestPrintNet(c *C) {
	ctx, cancel := context.WithCancel(context.Background())

	w := bytes.NewBufferString("")
	n := cpn.NewPN()
	n.P("pin",
		cpn.WithContext(ctx),
		cpn.WithPlace(memory.NewBlock()),
	)
	n.T("t1", cpn.WithFunction(transition.First))
	n.P("pout",
		cpn.WithContext(ctx),
		cpn.WithPlace(io.NewWriter(io.WriterOption(w))),
	)
	n.
		PT("pin", "t1").
		TP("t1", "pout")
	go func() {
		for i := 0; i < 10; i += 1 {
			n.P("pin").In() <- cpn.NewM(i)
		}
		cancel()
	}()
	n.RunSync()
	c.Assert(w.String(), Equals, "0\n1\n2\n3\n4\n5\n6\n7\n8\n9\n")
}

func (s *PNSuite) TestMultipleValues(c *C) {
	var n = cpn.NewPN()
	n.P("p1",
		cpn.WithContext(context.Background()),
		cpn.WithPlace(memory.NewBlock()),
	)
	n.T("t1", cpn.WithFunction(transition.First))
	n.P("p2",
		cpn.WithContext(context.Background()),
		cpn.WithPlace(place.NewPass(place.PassFuncOption(
			func(ctx context.Context, m *cpn.M) {
				m.SetValue("overwritten value from p2")
				m.SetValue("value from p2")
			},
		))),
	)
	n.T("t2", cpn.WithFunction(transition.First))
	n.P("p3",
		cpn.WithContext(context.Background()),
		cpn.WithPlace(place.NewPass(place.PassFuncOption(
			func(ctx context.Context, m *cpn.M) {
				m.SetValue("overwritten value from p3")
				m.SetValue("value from p3")
			},
		))),
	)
	n.
		PT("p1", "t1").
		TP("t1", "p2").
		PT("p2", "t2").
		TP("t2", "p3")

	var m = cpn.NewM("initial value")
	go func() {
		n.P("p1").In() <- m
		n.P("p1").Close()
	}()
	n.RunSync()

	// The backward compatibility check
	c.Assert(m.Value(), NotNil)
	c.Assert(m.Value().(string), Equals, "value from p3")

	// Verify that the mark values from different places are available
	c.Assert(m.ValueByPlace("not found", 0), IsNil)
	c.Assert(m.ValueByPlace("p1", 0), NotNil)
	c.Assert(m.ValueByPlace("p1", 0).(string), Equals, "initial value")
	c.Assert(m.ValueByPlace("p2", 0), NotNil)
	c.Assert(m.ValueByPlace("p2", 0).(string), Equals, "value from p2")
	c.Assert(m.ValueByPlace("p3", 0), NotNil)
	c.Assert(m.ValueByPlace("p3", 0).(string), Equals, "value from p3")
}
