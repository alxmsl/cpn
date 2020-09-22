package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/alxmsl/cpn"
	"github.com/alxmsl/cpn/place/memory"
	"github.com/alxmsl/cpn/place/redis"
	"github.com/alxmsl/cpn/transition"
	"github.com/mediocregopher/radix"
)

func main() {
	pool, err := radix.NewPool("tcp", "127.0.0.1:6379", 1,
		radix.PoolPingInterval(time.Second),
	)
	if err != nil {
		log.Println("error on create radix pool", err)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	n := cpn.NewPN()
	n.P("pin", cpn.WithContext(ctx),
		cpn.WithPlace(memory.NewBlock()),
	)
	n.T("t", cpn.WithFunction(transition.First))
	n.P("queue", cpn.WithContext(context.Background()),
		cpn.WithPlaceBuilder(redis.NewPush,
			redis.PoolOption(pool),
			redis.KeyOption("my_queue"),
			redis.MarshallerOption(redis.JsonMarshal),
		),
	)
	n.
		PT("pin", "t").
		TP("t", "queue")
	go func() {
		const numberOfTokens = 10
		for i := 0; i < numberOfTokens; i += 1 {
			n.P("pin").In() <- cpn.NewM(i)
		}
		fmt.Printf("wrote tokens: %d\n", numberOfTokens)
		cancel()
	}()
	n.RunSync()

	var v string
	err = pool.Do(radix.Cmd(&v, "LLEN", "my_queue"))
	if err != nil {
		log.Println("error on get queue length", err)
		return
	}
	fmt.Printf("queue lenth: %s\n", v)
}
