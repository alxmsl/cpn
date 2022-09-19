package test

import (
	"context"

	. "gopkg.in/check.v1"

	"github.com/alxmsl/cpn"
	"github.com/alxmsl/cpn/strategies"
	"github.com/alxmsl/cpn/transition"
)

type StrategiesSuite struct{}

var _ = Suite(&StrategiesSuite{})

func (s *StrategiesSuite) TestForkStrategy(c *C) {
	var n = cpn.NewPN()
	n.P("pin",
		cpn.WithContext(context.Background()),
		cpn.WithStrategy(strategies.NewFork(strategies.ForkFuncOption(
			func(ctx context.Context, m *cpn.M, ch chan<- *cpn.M) {
				m.SetValue("value from pin")
				ch <- m
			},
		))),
	)

	// The first branch
	n.T("t1.1", cpn.WithTransformation(transition.First))
	n.P("p1",
		cpn.WithContext(context.Background()),
		cpn.WithStrategy(strategies.NewPass(strategies.PassFuncOption(
			func(ctx context.Context, m *cpn.M) *cpn.M {
				m.SetValue("value from p1")
				return m
			},
		))),
	)
	n.T("t1.2", cpn.WithTransformation(transition.First))
	n.P("pout1",
		cpn.WithContext(context.Background()),
		cpn.WithStrategy(strategies.NewPass(strategies.PassFuncOption(
			func(ctx context.Context, m *cpn.M) *cpn.M {
				m.SetValue("value from pout1")
				return m
			},
		))),
		cpn.WithKeep(true),
	)

	// The second branch
	n.T("t2.1", cpn.WithTransformation(transition.First))
	n.P("p2",
		cpn.WithContext(context.Background()),
		cpn.WithStrategy(strategies.NewPass(strategies.PassFuncOption(
			func(ctx context.Context, m *cpn.M) *cpn.M {
				m.SetValue("value from p2")
				return m
			},
		))),
	)
	n.T("t2.2", cpn.WithTransformation(transition.First))
	n.P("pout2",
		cpn.WithContext(context.Background()),
		cpn.WithStrategy(strategies.NewPass(strategies.PassFuncOption(
			func(ctx context.Context, m *cpn.M) *cpn.M {
				m.SetValue("value from pout2")
				return m
			},
		))),
		cpn.WithKeep(true),
	)

	n.
		PT("pin", "t1.1").
		TP("t1.1", "p1").
		PT("p1", "t1.2").
		TP("t1.2", "pout1").
		PT("pin", "t2.1").
		TP("t2.1", "p2").
		PT("p2", "t2.2").
		TP("t2.2", "pout2").
		Run()

	// Writes several token to the PN concurrently. We use bigger number to be sure tokens are passed through several
	// branches. With small number tokens may be passed through one branch
	const loops = 10
	for i := 0; i < loops; i += 1 {
		n.P("pin").In() <- cpn.NewM("initial value")
	}

	var passed1, passed2 bool
	// Checks tokens are passed the PN in general
	for i := 0; i < loops; i += 1 {
		select {
		case m := <-n.P("pout1").Out():
			c.Assert(m.Word(), HasLen, 2)
			c.Assert(m.Word()[0], Equals, "t1.1")
			c.Assert(m.Word()[1], Equals, "t1.2")
			passed1 = true
		case m := <-n.P("pout2").Out():
			c.Assert(m.Word(), HasLen, 2)
			c.Assert(m.Word()[0], Equals, "t2.1")
			c.Assert(m.Word()[1], Equals, "t2.2")
			passed2 = true
		}
	}

	// Awaits tokens are passed on both branches in the PN
	c.Assert(passed1, Equals, true)
	c.Assert(passed2, Equals, true)
}

func (s *StrategiesSuite) TestJoinStrategy(c *C) {
	var (
		ctx, cancel = context.WithCancel(context.Background())
		n           = cpn.NewPN()
	)

	// The first branch
	n.P("pin1",
		cpn.WithContext(ctx),
		cpn.WithStrategy(strategies.NewPass(strategies.PassFuncOption(
			func(ctx context.Context, m *cpn.M) *cpn.M {
				m.SetValue("value from pin1")
				return m
			},
		))),
	)
	n.T("t1.1", cpn.WithTransformation(transition.First))

	// The second branch
	n.P("pin2",
		cpn.WithContext(ctx),
		cpn.WithStrategy(strategies.NewPass(strategies.PassFuncOption(
			func(ctx context.Context, m *cpn.M) *cpn.M {
				m.SetValue("value from pin2")
				return m
			},
		))),
	)
	n.T("t2.1", cpn.WithTransformation(transition.First))

	n.P("pout",
		cpn.WithContext(context.Background()),
		cpn.WithStrategy(strategies.NewJoin(strategies.JoinFuncOption(
			func(ctx context.Context, ch <-chan *cpn.M) *cpn.M {
				var counter int
				for range ch {
					counter += 1
				}
				return cpn.NewM(counter)
			},
		))),
		cpn.WithKeep(true),
	)

	n.
		PT("pin1", "t1.1").
		PT("pin2", "t2.1").
		TP("t1.1", "pout").
		TP("t2.1", "pout").
		Run()

	// Writes several token to the PN concurrently
	n.P("pin1").Send(cpn.NewM("initial value 1"))
	n.P("pin2").Send(cpn.NewM("initial value 2"))
	cancel()

	// Expects an aggregated token
	m := <-n.P("pout").Out()
	c.Assert(m.Value(), NotNil)
	v, ok := m.Value().(int)
	c.Assert(ok, Equals, true)
	c.Assert(v, Equals, 2)
}

func (s *StrategiesSuite) TestPassStrategy(c *C) {
	var n = cpn.NewPN()
	n.P("pin",
		cpn.WithContext(context.Background()),
		cpn.WithStrategy(strategies.NewPass(strategies.PassFuncOption(
			func(ctx context.Context, m *cpn.M) *cpn.M {
				m.SetValue("value from pin")
				return m
			},
		))),
	)
	n.T("t1", cpn.WithTransformation(transition.First))
	n.P("p",
		cpn.WithContext(context.Background()),
		cpn.WithStrategy(strategies.NewPass(strategies.PassFuncOption(
			func(ctx context.Context, m *cpn.M) *cpn.M {
				m.SetValue("value from p2")
				return m
			},
		))),
	)
	n.T("t2", cpn.WithTransformation(transition.First))
	n.P("pout",
		cpn.WithContext(context.Background()),
		cpn.WithStrategy(strategies.NewPass(strategies.PassFuncOption(
			func(ctx context.Context, m *cpn.M) *cpn.M {
				m.SetValue("value from p3")
				return m
			},
		))),
		cpn.WithKeep(true),
	)
	n.
		PT("pin", "t1").
		TP("t1", "p").
		PT("p", "t2").
		TP("t2", "pout").
		Run()

	n.P("pin").In() <- cpn.NewM("initial value")
	var m = <-n.P("pout").Out()
	c.Assert(m.Word(), HasLen, 2)
	c.Assert(m.Word()[0], Equals, "t1")
	c.Assert(m.Word()[1], Equals, "t2")
}
