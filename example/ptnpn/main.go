package main

import (
	"context"
	"fmt"

	"github.com/alxmsl/cpn"
	"github.com/alxmsl/cpn/place/memory"
	"github.com/alxmsl/cpn/transition"
)

const places = 3

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	n := cpn.NewPN()
	n.P("pin",
		cpn.WithContext(ctx),
		cpn.WithPlace(memory.NewBlock()),
		//cpn.IsInitial(),
	)
	n.Tn(places, "t_", cpn.WithFunction(transition.First))
	n.Pn(places, "p_",
		cpn.WithContext(context.Background()),
		cpn.WithPlaceBuilder(memory.NewBlock),
		//cpn.IsFinal(),
	)
	n.PTn(places, "pin", "t_").TnPn(places, "t_", "p_")

	go func() {
		for i := 0; i < 10; i += 1 {
			n.P("pin").In() <- cpn.NewM(i)
		}
		cancel()
	}()
	for i := 0; i < places; i += 1 {
		go func(i int) {
			for m := range n.P(fmt.Sprintf("p_:%d", i)).Out() {
				fmt.Println(m)
			}
		}(i)
	}
	n.RunSync()

	k, m := n.Size()
	fmt.Printf("pn size: %dx%d\n", k, m)
}
