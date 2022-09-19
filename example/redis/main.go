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
	fmt.Println("push to the queue:")
	n1 := newPushPN()
	n1.RunSync()

	printQueueLength()

	fmt.Println("pop from the queue:")
	n2 := newPopPN()
	n2.RunSync()
}

type MyType int

func (i *MyType) String() string {
	return fmt.Sprintf("%d\n", int(*i))
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
		cpn.WithStrategy(memory.NewBlock()),
	)
	n.T("t", cpn.WithTransformation(transition.First))
	n.P("queue", cpn.WithContext(context.Background()),
		cpn.WithStrategyBuilder(redis.NewPush,
			redis.PoolOption(pool),
			redis.KeyOption("my_queue"),
			redis.MarshallerOption(redis.JsonMarshal),
		),
	)
	n.P("print",
		cpn.WithContext(context.Background()),
		cpn.WithStrategy(io.NewWriter(io.WriterOption(os.Stdout))),
	)
	n.
		PT("pin", "t").
		TP("t", "queue").
		TP("t", "print")
	go func() {
		for i := 0; i < numberOfTokens; i += 1 {
			n.P("pin").In() <- cpn.NewM(i)
		}
		n.P("pin").Close()
	}()
	return n
}

func newPopPN() *cpn.PN {
	n := cpn.NewPN()
	n.P("queue", cpn.WithContext(context.Background()),
		cpn.WithStrategyBuilder(redis.NewPop,
			redis.PoolOption(pool),
			redis.KeyOption("my_queue"),
			redis.TypeOption(reflect.TypeOf(MyType(0))),
			redis.UnmarshallerOption(redis.JsonUnmarshal),
		),
	)
	n.T("t", cpn.WithTransformation(transition.First))
	n.P("print",
		cpn.WithContext(context.Background()),
		cpn.WithStrategy(io.NewWriter(io.WriterOption(os.Stdout))),
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
