package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/alxmsl/rtpn"
	"github.com/alxmsl/rtpn/place"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	var (
		pin = rtpn.NewP("pin").
			SetOptions(rtpn.WithContextOption(ctx)).
			SetOptions(rtpn.WithPlaceOption(place.NewBlock()))
		pout = rtpn.NewP("pout").
			SetOptions(rtpn.WithContextOption(ctx)).
			SetOptions(rtpn.WithPlaceOption(place.NewBlock()))
		t1 = rtpn.NewT("t1")
	)

	n := rtpn.NewPN()
	n.PT(pin, t1)
	n.TP(t1, pout)
	n.Run()

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 10; i += 1 {
			pin.WriteCh() <- rtpn.NewM(i)
		}
		cancel()
	}()
	go func() {
		defer wg.Done()
		for {
			m, ok := pout.Read()
			if !ok {
				return
			}
			fmt.Println(m)
		}
	}()
	wg.Wait()

	k, m := n.Size()
	fmt.Printf("pn size: %dx%d\n", k, m)
}
