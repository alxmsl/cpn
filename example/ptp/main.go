package main

import (
	"context"
	"fmt"

	"github.com/alxmsl/cpn"
	"github.com/alxmsl/cpn/place"
	"github.com/alxmsl/cpn/transition"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	n := cpn.NewPN()
	n.P("pin", cpn.WithContext(ctx), cpn.WithPlace(place.NewBlock()))
	n.T("t1", cpn.WithFunction(transition.First))
	n.P("pout", cpn.WithContext(ctx), cpn.WithPlace(place.NewBlock()), cpn.IsTermination())

	n.PT("pin", "t1").TP("t1", "pout").Run()

	go func() {
		for i := 0; i < 10; i += 1 {
			n.P("pin").In() <- cpn.NewM(i)
		}
		cancel()
	}()
	for m := range n.P("pout").Out() {
		fmt.Println(m)
	}

	k, m := n.Size()
	fmt.Printf("pn size: %dx%d\n", k, m)
}
