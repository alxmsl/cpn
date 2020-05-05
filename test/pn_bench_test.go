package test

import (
	"context"
	"testing"

	"github.com/alxmsl/rtpn"
	"github.com/alxmsl/rtpn/place"
	"github.com/alxmsl/rtpn/transition"
)

func BenchmarkBlockPTP(b *testing.B) {
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

	mm := make([]*rtpn.M, b.N)
	for i := 0; i < b.N; i += 1 {
		mm[i] = rtpn.NewM(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i += 1 {
		pin.WriteCh() <- mm[i]
		pout.Read()
	}
}

func BenchmarkQueuePTP(b *testing.B) {
	pin := rtpn.NewP("pin").
		SetOptions(rtpn.WithContext(context.Background())).
		SetOptions(rtpn.WithPlace(place.NewQueue(100)))
	pout := rtpn.NewP("pout").
		SetOptions(rtpn.WithContext(context.Background())).
		SetOptions(rtpn.WithPlace(place.NewQueue(100)))
	t1 := rtpn.NewT("t1").
		SetOptions(rtpn.WithFunction(transition.First))

	n := rtpn.NewPN()
	n.PT(pin, t1)
	n.TP(t1, pout)
	n.Run()

	mm := make([]*rtpn.M, b.N)
	for i := 0; i < b.N; i += 1 {
		mm[i] = rtpn.NewM(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i += 1 {
		pin.WriteCh() <- mm[i]
		pout.Read()
	}
}

func BenchmarkPPTP(b *testing.B) {
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

	mm := make([]*rtpn.M, b.N)
	for i := 0; i < b.N; i += 1 {
		mm[i] = rtpn.NewM(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i += 1 {
		p1.WriteCh() <- mm[i]
		p2.WriteCh() <- mm[i]
		pout.Read()
	}
}

func BenchmarkPPTT(b *testing.B) {
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

	mm := make([]*rtpn.M, b.N)
	for i := 0; i < b.N; i += 1 {
		mm[i] = rtpn.NewM(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i += 1 {
		p1.WriteCh() <- mm[i]
		p2.WriteCh() <- mm[i]
		pout.Read()
	}
}
