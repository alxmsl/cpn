package test

import (
	"context"
	"testing"

	"github.com/alxmsl/cpn"
	"github.com/alxmsl/cpn/place"
	"github.com/alxmsl/cpn/place/memory"
	"github.com/alxmsl/cpn/transition"
)

func BenchmarkBlockPTP(b *testing.B) {
	n := cpn.NewPN()
	n.P("pin", cpn.WithContext(context.Background()), cpn.WithPlace(place.NewBlock()))
	n.T("t1", cpn.WithFunction(transition.First))
	n.P("pout", cpn.WithContext(context.Background()), cpn.WithPlace(place.NewBlock()))

	n.PT("pin", "t1").TP("t1", "pout").Run()

	mm := make([]*cpn.M, b.N)
	for i := 0; i < b.N; i += 1 {
		mm[i] = cpn.NewM(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i += 1 {
		n.P("pin").In() <- mm[i]
		n.P("pout").Read()
	}
}

func BenchmarkBlockPTPTP(b *testing.B) {
	n := cpn.NewPN()
	n.P("pin", cpn.WithContext(context.Background()), cpn.WithPlace(place.NewBlock()))
	n.P("p1", cpn.WithContext(context.Background()), cpn.WithPlace(place.NewBlock()))
	n.P("pout", cpn.WithContext(context.Background()), cpn.WithPlace(place.NewBlock()))
	n.T("t1", cpn.WithFunction(transition.First))
	n.T("t2", cpn.WithFunction(transition.First))

	n.PT("pin", "t1").TP("t1", "p1").PT("p1", "t2").TP("t2", "pout").Run()

	mm := make([]*cpn.M, b.N)
	for i := 0; i < b.N; i += 1 {
		mm[i] = cpn.NewM(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i += 1 {
		n.P("pin").In() <- mm[i]
		n.P("pout").Read()
	}
}

func BenchmarkBlockPTPTPTP(b *testing.B) {
	n := cpn.NewPN()
	n.P("pin", cpn.WithContext(context.Background()), cpn.WithPlace(place.NewBlock()))
	n.P("p1", cpn.WithContext(context.Background()), cpn.WithPlace(place.NewBlock()))
	n.P("p2", cpn.WithContext(context.Background()), cpn.WithPlace(place.NewBlock()))
	n.P("pout", cpn.WithContext(context.Background()), cpn.WithPlace(place.NewBlock()))
	n.T("t1", cpn.WithFunction(transition.First))
	n.T("t2", cpn.WithFunction(transition.First))
	n.T("t3", cpn.WithFunction(transition.First))

	n.PT("pin", "t1").TP("t1", "p1").PT("p1", "t2").TP("t2", "p2").PT("p2", "t3").
		TP("t3", "pout").Run()

	mm := make([]*cpn.M, b.N)
	for i := 0; i < b.N; i += 1 {
		mm[i] = cpn.NewM(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i += 1 {
		n.P("pin").In() <- mm[i]
		n.P("pout").Read()
	}
}

func BenchmarkQueuePTP(b *testing.B) {
	n := cpn.NewPN()
	n.P("pin", cpn.WithContext(context.Background()), cpn.WithPlace(memory.NewQueue(100)))
	n.T("t1", cpn.WithFunction(transition.First))
	n.P("pout", cpn.WithContext(context.Background()), cpn.WithPlace(memory.NewQueue(100)))

	n.PT("pin", "t1").TP("t1", "pout").Run()

	mm := make([]*cpn.M, b.N)
	for i := 0; i < b.N; i += 1 {
		mm[i] = cpn.NewM(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i += 1 {
		n.P("pin").In() <- mm[i]
		n.P("pout").Read()
	}
}

func BenchmarkPPTP(b *testing.B) {
	n := cpn.NewPN()
	n.P("p1", cpn.WithContext(context.Background()), cpn.WithPlace(place.NewBlock()))
	n.P("p2", cpn.WithContext(context.Background()), cpn.WithPlace(place.NewBlock()))
	n.T("t1", cpn.WithFunction(transition.First))
	n.P("pout", cpn.WithContext(context.Background()), cpn.WithPlace(place.NewBlock()))

	n.PT("p1", "t1").PT("p2", "t1").TP("t1", "pout").Run()

	mm := make([]*cpn.M, b.N)
	for i := 0; i < b.N; i += 1 {
		mm[i] = cpn.NewM(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i += 1 {
		n.P("p1").In() <- mm[i]
		n.P("p2").In() <- mm[i]
		n.P("pout").Read()
	}
}

func BenchmarkPPTTP(b *testing.B) {
	n := cpn.NewPN()
	n.P("p1", cpn.WithContext(context.Background()), cpn.WithPlace(place.NewBlock()))
	n.P("p2", cpn.WithContext(context.Background()), cpn.WithPlace(place.NewBlock()))
	n.T("t1", cpn.WithFunction(transition.First))
	n.T("t2", cpn.WithFunction(transition.First))
	n.P("pout", cpn.WithContext(context.Background()), cpn.WithPlace(place.NewBlock()))

	n.PT("p1", "t1").PT("p2", "t1").PT("p1", "t2").PT("p2", "t2").
		TP("t1", "pout").TP("t2", "pout").Run()

	mm := make([]*cpn.M, b.N)
	for i := 0; i < b.N; i += 1 {
		mm[i] = cpn.NewM(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i += 1 {
		n.P("p1").In() <- mm[i]
		n.P("p2").In() <- mm[i]
		n.P("pout").Read()
	}
}
