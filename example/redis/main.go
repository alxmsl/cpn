package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"reflect"
	"time"

	"github.com/alxmsl/cpn"
	"github.com/alxmsl/cpn/place/io"
	"github.com/alxmsl/cpn/place/memory"
	"github.com/alxmsl/cpn/place/redis"
	"github.com/alxmsl/cpn/transition"
	"github.com/mediocregopher/radix/v3"
)

const numberOfTokens = 10

var pool *radix.Pool

func init() {
	var err error
	pool, err = radix.NewPool("tcp", "127.0.0.1:6379", 1,
		radix.PoolPingInterval(time.Second),
	)
	if err != nil {
		log.Println("error on create radix pool", err)
		os.Exit(1)
	}
}

func main() {
	n1 := newPushPN()
	n1.RunSync()

	printQueueLength()

	n2 := newPopPN()
	n2.RunSync()
}

func printQueueLength() {
	var v string
	err := pool.Do(radix.Cmd(&v, "LLEN", "my_queue"))
	if err != nil {
		log.Println("error on get queue length", err)
		return
	}
	fmt.Printf("queue lenth: %s\n", v)
}

func newPushPN() *cpn.PN {
	n := cpn.NewPN()
	n.P("pin", cpn.WithContext(context.Background()),
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
		for i := 0; i < numberOfTokens; i += 1 {
			n.P("pin").In() <- cpn.NewM(i)
		}
		fmt.Printf("wrote tokens: %d\n", numberOfTokens)
		n.P("pin").Close()
	}()
	return n
}

func newPopPN() *cpn.PN {
	n := cpn.NewPN()
	n.P("queue", cpn.WithContext(context.Background()),
		cpn.WithPlaceBuilder(redis.NewPop,
			redis.PoolOption(pool),
			redis.KeyOption("my_queue"),
			redis.TypeOption(reflect.TypeOf(0)),
			redis.UnmarshallerOption(redis.JsonUnmarshal),
		),
	)
	n.T("t", cpn.WithFunction(transition.First))
	n.P("print",
		cpn.WithContext(context.Background()),
		cpn.WithPlace(io.NewWriter(io.WriterOption(os.Stdout))),
	)
	n.
		PT("queue", "t").
		TP("t", "print")
	go func() {
		for i := 0; i < numberOfTokens; i += 1 {
			n.P("queue").In() <- cpn.NewM(nil)
		}
		n.P("queue").Close()
	}()
	return n
}
