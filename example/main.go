package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/alxmsl/pn"
	"github.com/alxmsl/pn/place"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	var (
		pin = pn.NewP("pin").
			SetOptions(pn.WithContextOption(ctx)).
			SetOptions(pn.WithPlaceOption(place.NewBlock()))
		pout = pn.NewP("pout").
			SetOptions(pn.WithContextOption(ctx)).
			SetOptions(pn.WithPlaceOption(place.NewBlock()))
		t1 = pn.NewT("t1")
	)

	n := pn.NewPN()
	n.PT(pin, t1)
	n.TP(t1, pout)
	n.Run()

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 10; i += 1 {
			pin.WriteCh() <- pn.NewM(i)
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
