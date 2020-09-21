package test

import (
	"log"

	. "gopkg.in/check.v1"

	"bytes"
	"context"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/alxmsl/cpn"
	"github.com/alxmsl/cpn/place/memory"
	"github.com/alxmsl/cpn/place/out"
	"github.com/alxmsl/cpn/transition"
)

func Test(t *testing.T) {
	TestingT(t)
}

type PNSuite struct{}

var _ = Suite(&PNSuite{})

func (s *PNSuite) TestPTP(c *C) {
	c.Skip("1")

	ctx, cancel := context.WithCancel(context.Background())

	n := cpn.NewPN()
	n.P("pin",
		cpn.WithContext(ctx),
		cpn.WithPlace(memory.NewBlock()),
		//cpn.IsInitial(),
	)
	n.T("t1", cpn.WithFunction(transition.First))
	n.P("pout",
		cpn.WithContext(context.Background()),
		cpn.WithPlace(memory.NewBlock()),
		//cpn.IsFinal(),
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

		i += 1
	}
}

func (s *PNSuite) TestPTTP(c *C) {
	c.Skip("1")

	ctx, cancel := context.WithCancel(context.Background())

	n := cpn.NewPN()
	n.P("pin",
		cpn.WithContext(ctx),
		cpn.WithPlace(memory.NewBlock()),
		//cpn.IsInitial(),
	)
	n.T("t1", cpn.WithFunction(transition.First))
	n.T("t2", cpn.WithFunction(transition.First))
	n.P("pout",
		cpn.WithContext(ctx),
		cpn.WithPlace(memory.NewBlock()),
		//cpn.IsFinal(),
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

		count += 1
	}
	c.Assert(count, Equals, 1000)
}

func (s *PNSuite) TestPTPP(c *C) {
	c.Skip("1")

	ctx, cancel := context.WithCancel(context.Background())

	n := cpn.NewPN()
	n.P("pin",
		cpn.WithContext(ctx),
		cpn.WithPlace(memory.NewBlock()),
		//cpn.IsInitial(),
	)
	n.T("t", cpn.WithFunction(transition.First))
	n.P("pout1",
		cpn.WithContext(ctx),
		cpn.WithPlace(memory.NewBlock()),
		//cpn.IsFinal(),
	)
	n.P("pout2",
		cpn.WithContext(ctx),
		cpn.WithPlace(memory.NewBlock()),
		//cpn.IsFinal(),
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
		for range n.P("pout1").Out() {
			atomic.AddInt64(&count, 1)
		}
	}()
	go func() {
		defer wg.Done()
		for range n.P("pout2").Out() {
			atomic.AddInt64(&count, 1)
		}
	}()
	wg.Wait()
	c.Assert(count, Equals, int64(2000))
}

func (s *PNSuite) TestPPTP(c *C) {
	c.Skip("1")

	ctx, cancel := context.WithCancel(context.Background())

	n := cpn.NewPN()
	n.P("p1",
		cpn.WithContext(ctx),
		cpn.WithPlace(memory.NewBlock()),
		//cpn.IsInitial(),
	)
	n.P("p2",
		cpn.WithContext(ctx),
		cpn.WithPlace(memory.NewBlock()),
		//cpn.IsInitial(),
	)
	n.T("t1", cpn.WithFunction(transition.First))
	n.P("pout",
		cpn.WithContext(ctx),
		cpn.WithPlace(memory.NewBlock()),
		//cpn.IsFinal(),
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
		c.Assert(m.Value().(int), Equals, i)
		c.Assert(m.Path(), HasLen, 3)
		c.Assert(m.Path()[0].N == "p1" || m.Path()[0].N == "p2", Equals, true)
		c.Assert(m.Path()[1].N, Equals, "t1")
		c.Assert(m.Path()[2].N, Equals, "pout")

		c.Assert(m.History(), HasLen, 4)
		c.Assert(m.History()[0].N, Equals, "")
		c.Assert(m.History()[1], Equals, m.Path()[0])
		c.Assert(m.History()[2], Equals, m.Path()[1])
		c.Assert(m.History()[3], Equals, m.Path()[2])

		i += 1
	}
}

func (s *PNSuite) TestPPTTP(c *C) {
	c.Skip("1")

	ctx, cancel := context.WithCancel(context.Background())

	n := cpn.NewPN()
	n.P("p1",
		cpn.WithContext(ctx),
		cpn.WithPlace(memory.NewBlock()),
		//cpn.IsInitial(),
	)
	n.T("t1", cpn.WithFunction(transition.First))
	n.P("p2",
		cpn.WithContext(ctx),
		cpn.WithPlace(memory.NewBlock()),
		//cpn.IsInitial(),
	)
	n.T("t2", cpn.WithFunction(transition.First))
	n.P("pout",
		cpn.WithContext(ctx),
		cpn.WithPlace(memory.NewBlock()),
		//cpn.IsFinal(),
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
		c.Assert(m.Path()[0].N == "p1" || m.Path()[0].N == "p2", Equals, true)
		c.Assert(m.Path()[1].N == "t1" || m.Path()[1].N == "t2", Equals, true)
		c.Assert(m.Path()[2].N, Equals, "pout")

		c.Assert(m.History(), HasLen, 4)
		c.Assert(m.History()[0].N, Equals, "")
		c.Assert(m.History()[1], Equals, m.Path()[0])
		c.Assert(m.History()[2], Equals, m.Path()[1])
		c.Assert(m.History()[3], Equals, m.Path()[2])

		count += 1
	}
	c.Assert(count, Equals, 1000)
}

func (s *PNSuite) TestPrintNet(c *C) {
	//c.Skip("1")

	ctx, cancel := context.WithCancel(context.Background())

	w := bytes.NewBufferString("")
	n := cpn.NewPN()
	n.P("pin",
		cpn.WithContext(ctx),
		cpn.WithPlace(memory.NewBlock()),
		//cpn.IsInitial(),
	)
	n.T("t1", cpn.WithFunction(transition.First))
	n.P("pout",
		cpn.WithContext(context.Background()),
		cpn.WithPlaceBuilder(out.NewPrint,
			out.WriterOption(w),
		),
		//cpn.IsFinal(),
	)
	//n.T("null", cpn.WithFunction(transition.Null))

	n.
		PT("pin", "t1").
		TP("t1", "pout")
		//PT("pout", "null")
	go func() {
		for i := 0; i < 10; i += 1 {
			log.Println(i)
			n.P("pin").In() <- cpn.NewM(i)
		}
		cancel()
	}()
	n.RunSync()
	c.Assert(w.String(), Equals, "0\n1\n2\n3\n4\n5\n6\n7\n8\n9\n")
}
