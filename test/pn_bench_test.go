package test

import (
	"context"
	"testing"

	"github.com/alxmsl/pn"
	"github.com/alxmsl/pn/place"
)

func BenchmarkBlockPTP(b *testing.B) {
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

	mm := make([]*pn.M, b.N)
	for i := 0; i < b.N; i += 1 {
		mm[i] = pn.NewM(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i += 1 {
		pin.WriteCh() <- mm[i]
		pout.Read()
	}
}

func BenchmarkQueuePTP(b *testing.B) {
	pin := pn.NewP("pin").
		SetOptions(pn.WithContextOption(context.Background())).
		SetOptions(pn.WithPlaceOption(place.NewQueue(100)))
	pout := pn.NewP("pout").
		SetOptions(pn.WithContextOption(context.Background())).
		SetOptions(pn.WithPlaceOption(place.NewQueue(100)))
	t1 := pn.NewT("t1")

	n := pn.NewPN()
	n.PT(pin, t1)
	n.TP(t1, pout)
	n.Run()

	mm := make([]*pn.M, b.N)
	for i := 0; i < b.N; i += 1 {
		mm[i] = pn.NewM(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i += 1 {
		pin.WriteCh() <- mm[i]
		pout.Read()
	}
}

func BenchmarkPPTP(b *testing.B) {
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

	mm := make([]*pn.M, b.N)
	for i := 0; i < b.N; i += 1 {
		mm[i] = pn.NewM(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i += 1 {
		p1.WriteCh() <- mm[i]
		p2.WriteCh() <- mm[i]
		pout.Read()
	}
}

func BenchmarkPPTT(b *testing.B) {
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

	mm := make([]*pn.M, b.N)
	for i := 0; i < b.N; i += 1 {
		mm[i] = pn.NewM(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i += 1 {
		p1.WriteCh() <- mm[i]
		p2.WriteCh() <- mm[i]
		pout.Read()
	}
}
