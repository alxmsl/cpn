package test

import (
	"context"
	"testing"

	"github.com/alxmsl/cpn"
	"github.com/alxmsl/cpn/place/memory"
	"github.com/alxmsl/cpn/transition"
)

func BenchmarkBlockPTP(b *testing.B) {
	n := cpn.NewPN()
	n.P("pin",
		cpn.WithContext(context.Background()),
		cpn.WithPlace(memory.NewBlock()),
	)
	n.T("t1", cpn.WithFunction(transition.First))
	n.P("pout",
		cpn.WithContext(context.Background()),
		cpn.WithPlace(memory.NewBlock()),
	)

	mm := make([]*cpn.M, b.N)
	for i := 0; i < b.N; i += 1 {
		mm[i] = cpn.NewM(i)
	}

	n.
		PT("pin", "t1").
		TP("t1", "pout").
		Run()
	b.ResetTimer()
	for i := 0; i < b.N; i += 1 {
		n.P("pin").In() <- mm[i]
		<-n.P("pout").Out()
	}
}

func BenchmarkBlockPTPTP(b *testing.B) {
	n := cpn.NewPN()
	n.P("pin",
		cpn.WithContext(context.Background()),
		cpn.WithPlace(memory.NewBlock()),
	)
	n.T("t1", cpn.WithFunction(transition.First))
	n.P("p1",
		cpn.WithContext(context.Background()),
		cpn.WithPlace(memory.NewBlock()),
	)
	n.T("t2", cpn.WithFunction(transition.First))
	n.P("pout",
		cpn.WithContext(context.Background()),
		cpn.WithPlace(memory.NewBlock()),
	)

	mm := make([]*cpn.M, b.N)
	for i := 0; i < b.N; i += 1 {
		mm[i] = cpn.NewM(i)
	}

	n.
		PT("pin", "t1").
		TP("t1", "p1").
		PT("p1", "t2").
		TP("t2", "pout").
		Run()
	b.ResetTimer()
	for i := 0; i < b.N; i += 1 {
		n.P("pin").In() <- mm[i]
		<-n.P("pout").Out()
	}
}

func BenchmarkBlockPTPTPTP(b *testing.B) {
	n := cpn.NewPN()
	n.P("pin",
		cpn.WithContext(context.Background()),
		cpn.WithPlace(memory.NewBlock()),
	)
	n.T("t1", cpn.WithFunction(transition.First))
	n.P("p1",
		cpn.WithContext(context.Background()),
		cpn.WithPlace(memory.NewBlock()),
	)
	n.T("t2", cpn.WithFunction(transition.First))
	n.P("p2",
		cpn.WithContext(context.Background()),
		cpn.WithPlace(memory.NewBlock()),
	)
	n.T("t3", cpn.WithFunction(transition.First))
	n.P("pout",
		cpn.WithContext(context.Background()),
		cpn.WithPlace(memory.NewBlock()),
	)

	mm := make([]*cpn.M, b.N)
	for i := 0; i < b.N; i += 1 {
		mm[i] = cpn.NewM(i)
	}

	n.
		PT("pin", "t1").
		TP("t1", "p1").
		PT("p1", "t2").
		TP("t2", "p2").
		PT("p2", "t3").
		TP("t3", "pout").
		Run()
	b.ResetTimer()
	for i := 0; i < b.N; i += 1 {
		n.P("pin").In() <- mm[i]
		<-n.P("pout").Out()
	}
}

func BenchmarkQueuePTP(b *testing.B) {
	n := cpn.NewPN()
	n.P("pin",
		cpn.WithContext(context.Background()),
		cpn.WithPlace(memory.NewQueue(100)),
	)
	n.T("t1", cpn.WithFunction(transition.First))
	n.P("pout",
		cpn.WithContext(context.Background()),
		cpn.WithPlace(memory.NewQueue(100)),
	)

	mm := make([]*cpn.M, b.N)
	for i := 0; i < b.N; i += 1 {
		mm[i] = cpn.NewM(i)
	}

	n.
		PT("pin", "t1").
		TP("t1", "pout").
		Run()
	b.ResetTimer()
	for i := 0; i < b.N; i += 1 {
		n.P("pin").In() <- mm[i]
		<-n.P("pout").Out()
	}
}

func BenchmarkPTPP(b *testing.B) {
	n := cpn.NewPN()
	n.P("pin",
		cpn.WithContext(context.Background()),
		cpn.WithPlace(memory.NewBlock()),
	)
	n.T("t", cpn.WithFunction(transition.First))
	n.P("pout1",
		cpn.WithContext(context.Background()),
		cpn.WithPlace(memory.NewBlock()),
	)
	n.P("pout2",
		cpn.WithContext(context.Background()),
		cpn.WithPlace(memory.NewBlock()),
	)

	mm := make([]*cpn.M, b.N)
	for i := 0; i < b.N; i += 1 {
		mm[i] = cpn.NewM(i)
	}

	n.
		PT("pin", "t").
		TP("t", "pout1").
		TP("t", "pout2").
		Run()
	b.ResetTimer()
	for i := 0; i < b.N; i += 1 {
		n.P("pin").In() <- mm[i]
		<-n.P("pout1").Out()
		<-n.P("pout2").Out()
	}
}

func BenchmarkPPTP(b *testing.B) {
	n := cpn.NewPN()
	n.P("p1",
		cpn.WithContext(context.Background()),
		cpn.WithPlace(memory.NewBlock()),
	)
	n.P("p2",
		cpn.WithContext(context.Background()),
		cpn.WithPlace(memory.NewBlock()),
	)
	n.T("t1", cpn.WithFunction(transition.First))
	n.P("pout",
		cpn.WithContext(context.Background()),
		cpn.WithPlace(memory.NewBlock()),
	)

	mm := make([]*cpn.M, b.N)
	for i := 0; i < b.N; i += 1 {
		mm[i] = cpn.NewM(i)
	}

	n.
		PT("p1", "t1").
		PT("p2", "t1").
		TP("t1", "pout").
		Run()
	b.ResetTimer()
	for i := 0; i < b.N; i += 1 {
		n.P("p1").In() <- mm[i]
		n.P("p2").In() <- mm[i]
		<-n.P("pout").Out()
	}
}

func BenchmarkPPTTP(b *testing.B) {
	n := cpn.NewPN()
	n.P("p1",
		cpn.WithContext(context.Background()),
		cpn.WithPlace(memory.NewBlock()),
	)
	n.T("t1", cpn.WithFunction(transition.First))
	n.P("p2",
		cpn.WithContext(context.Background()),
		cpn.WithPlace(memory.NewBlock()),
	)
	n.T("t2", cpn.WithFunction(transition.First))
	n.P("pout",
		cpn.WithContext(context.Background()),
		cpn.WithPlace(memory.NewBlock()),
	)

	mm := make([]*cpn.M, b.N)
	for i := 0; i < b.N; i += 1 {
		mm[i] = cpn.NewM(i)
	}

	n.
		PT("p1", "t1").
		PT("p2", "t1").
		PT("p1", "t2").
		PT("p2", "t2").
		TP("t1", "pout").
		TP("t2", "pout").
		Run()
	b.ResetTimer()
	for i := 0; i < b.N; i += 1 {
		n.P("p1").In() <- mm[i]
		n.P("p2").In() <- mm[i]
		<-n.P("pout").Out()
	}
}
