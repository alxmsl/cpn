package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/alxmsl/cpn"
	"github.com/alxmsl/cpn/place"
	"github.com/alxmsl/cpn/transition"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	n := cpn.NewPN()
	n.P("pin", cpn.WithContext(ctx),
		cpn.WithPlace(place.NewBlock()),
	)
	n.T("t", cpn.WithFunction(transition.First))
	n.P("pout1",
		cpn.WithContext(context.Background()),
		cpn.WithPlace(place.NewBlock()),
	)
	n.P("pout2",
		cpn.WithContext(context.Background()),
		cpn.WithPlace(place.NewBlock()),
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

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		for m := range n.P("pout1").Out() {
			fmt.Println(m)
		}
	}()
	go func() {
		defer wg.Done()
		for m := range n.P("pout2").Out() {
			fmt.Println(m)
		}
	}()
	wg.Wait()

	k, m := n.Size()
	fmt.Printf("pn size: %dx%d\n", k, m)
}
