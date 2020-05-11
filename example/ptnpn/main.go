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
	n.P("pin", cpn.WithContext(ctx), cpn.WithPlace(place.NewBlock()))
	for i := 0; i < 3; i += 1 {
		nt := fmt.Sprintf("t%d", i)
		n.T(nt, cpn.WithFunction(transition.First))
		np := fmt.Sprintf("p%d", i)
		n.P(np, cpn.WithContext(ctx), cpn.WithPlace(place.NewBlock()), cpn.IsTermination())
		n.PT("pin", nt).TP(nt, np)
	}
	n.Run()

	go func() {
		for i := 0; i < 10; i += 1 {
			n.P("pin").In() <- cpn.NewM(i)
		}
		cancel()
	}()

	wg := sync.WaitGroup{}
	for i := 0; i < 3; i += 1 {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for m := range n.P(fmt.Sprintf("p%d", i)).Out() {
				fmt.Println(m)
			}
		}(i)
	}
	wg.Wait()

	k, m := n.Size()
	fmt.Printf("pn size: %dx%d\n", k, m)
}
