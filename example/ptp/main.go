package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/alxmsl/rtpn"
	"github.com/alxmsl/rtpn/place"
	"github.com/alxmsl/rtpn/transition"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	n := rtpn.NewPN()
	n.P("pin", rtpn.WithContext(ctx), rtpn.WithPlace(place.NewBlock()))
	n.T("t1", rtpn.WithFunction(transition.First))
	n.P("pout", rtpn.WithContext(ctx), rtpn.WithPlace(place.NewBlock()))

	n.PT("pin", "t1").TP("t1", "pout").Run()

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 10; i += 1 {
			n.P("pin").In() <- rtpn.NewM(i)
		}
		cancel()
	}()
	go func() {
		defer wg.Done()
		for {
			m, ok := n.P("pout").Read()
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
