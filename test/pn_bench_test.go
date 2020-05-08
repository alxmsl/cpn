package test

import (
	"context"
	"testing"

	"github.com/alxmsl/rtpn"
	"github.com/alxmsl/rtpn/place"
	"github.com/alxmsl/rtpn/place/memory"
	"github.com/alxmsl/rtpn/transition"
)

func BenchmarkBlockPTP(b *testing.B) {
	n := rtpn.NewPN()
	n.P("pin", rtpn.WithContext(context.Background()), rtpn.WithPlace(place.NewBlock()))
	n.T("t1", rtpn.WithFunction(transition.First))
	n.P("pout", rtpn.WithContext(context.Background()), rtpn.WithPlace(place.NewBlock()))

	n.PT("pin", "t1").TP("t1", "pout").Run()

	mm := make([]*rtpn.M, b.N)
	for i := 0; i < b.N; i += 1 {
		mm[i] = rtpn.NewM(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i += 1 {
		n.P("pin").In() <- mm[i]
		n.P("pout").Read()
	}
}

func BenchmarkBlockPTPTP(b *testing.B) {
	n := rtpn.NewPN()
	n.P("pin", rtpn.WithContext(context.Background()), rtpn.WithPlace(place.NewBlock()))
	n.P("p1", rtpn.WithContext(context.Background()), rtpn.WithPlace(place.NewBlock()))
	n.P("pout", rtpn.WithContext(context.Background()), rtpn.WithPlace(place.NewBlock()))
	n.T("t1", rtpn.WithFunction(transition.First))
	n.T("t2", rtpn.WithFunction(transition.First))

	n.PT("pin", "t1").TP("t1", "p1").PT("p1", "t2").TP("t2", "pout").Run()

	mm := make([]*rtpn.M, b.N)
	for i := 0; i < b.N; i += 1 {
		mm[i] = rtpn.NewM(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i += 1 {
		n.P("pin").In() <- mm[i]
		n.P("pout").Read()
	}
}

func BenchmarkBlockPTPTPTP(b *testing.B) {
	n := rtpn.NewPN()
	n.P("pin", rtpn.WithContext(context.Background()), rtpn.WithPlace(place.NewBlock()))
	n.P("p1", rtpn.WithContext(context.Background()), rtpn.WithPlace(place.NewBlock()))
	n.P("p2", rtpn.WithContext(context.Background()), rtpn.WithPlace(place.NewBlock()))
	n.P("pout", rtpn.WithContext(context.Background()), rtpn.WithPlace(place.NewBlock()))
	n.T("t1", rtpn.WithFunction(transition.First))
	n.T("t2", rtpn.WithFunction(transition.First))
	n.T("t3", rtpn.WithFunction(transition.First))

	n.PT("pin", "t1").TP("t1", "p1").PT("p1", "t2").TP("t2", "p2").PT("p2", "t3").
		TP("t3", "pout").Run()

	mm := make([]*rtpn.M, b.N)
	for i := 0; i < b.N; i += 1 {
		mm[i] = rtpn.NewM(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i += 1 {
		n.P("pin").In() <- mm[i]
		n.P("pout").Read()
	}
}

func BenchmarkQueuePTP(b *testing.B) {
	n := rtpn.NewPN()
	n.P("pin", rtpn.WithContext(context.Background()), rtpn.WithPlace(memory.NewQueue(100)))
	n.T("t1", rtpn.WithFunction(transition.First))
	n.P("pout", rtpn.WithContext(context.Background()), rtpn.WithPlace(memory.NewQueue(100)))

	n.PT("pin", "t1").TP("t1", "pout").Run()

	mm := make([]*rtpn.M, b.N)
	for i := 0; i < b.N; i += 1 {
		mm[i] = rtpn.NewM(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i += 1 {
		n.P("pin").In() <- mm[i]
		n.P("pout").Read()
	}
}

func BenchmarkPPTP(b *testing.B) {
	n := rtpn.NewPN()
	n.P("p1", rtpn.WithContext(context.Background()), rtpn.WithPlace(place.NewBlock()))
	n.P("p2", rtpn.WithContext(context.Background()), rtpn.WithPlace(place.NewBlock()))
	n.T("t1", rtpn.WithFunction(transition.First))
	n.P("pout", rtpn.WithContext(context.Background()), rtpn.WithPlace(place.NewBlock()))

	n.PT("p1", "t1").PT("p2", "t1").TP("t1", "pout").Run()

	mm := make([]*rtpn.M, b.N)
	for i := 0; i < b.N; i += 1 {
		mm[i] = rtpn.NewM(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i += 1 {
		n.P("p1").In() <- mm[i]
		n.P("p2").In() <- mm[i]
		n.P("pout").Read()
	}
}

func BenchmarkPPTTP(b *testing.B) {
	n := rtpn.NewPN()
	n.P("p1", rtpn.WithContext(context.Background()), rtpn.WithPlace(place.NewBlock()))
	n.P("p2", rtpn.WithContext(context.Background()), rtpn.WithPlace(place.NewBlock()))
	n.T("t1", rtpn.WithFunction(transition.First))
	n.T("t2", rtpn.WithFunction(transition.First))
	n.P("pout", rtpn.WithContext(context.Background()), rtpn.WithPlace(place.NewBlock()))

	n.PT("p1", "t1").PT("p2", "t1").PT("p1", "t2").PT("p2", "t2").
		TP("t1", "pout").TP("t2", "pout").Run()

	mm := make([]*rtpn.M, b.N)
	for i := 0; i < b.N; i += 1 {
		mm[i] = rtpn.NewM(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i += 1 {
		n.P("p1").In() <- mm[i]
		n.P("p2").In() <- mm[i]
		n.P("pout").Read()
	}
}
