package redis

import (
	"reflect"

	"github.com/alxmsl/cpn"
	"github.com/mediocregopher/radix/v3"
)

type Key interface {
	SetKey(string)
}

func KeyOption(k string) cpn.StrategyOption {
	return keyOption{k}
}

type keyOption struct {
	k string
}

func (o keyOption) Apply(p cpn.Strategy) {
	p.(Key).SetKey(o.k)
}

type MarshalFunc func(interface{}) (string, error)

func MarshallerOption(f MarshalFunc) cpn.StrategyOption {
	return marshallerOption{f}
}

type marshallerOption struct {
	f MarshalFunc
}

func (o marshallerOption) Apply(p cpn.Strategy) {
	p.(*Push).f = o.f
}

type Pool interface {
	SetPool(*radix.Pool)
}

func PoolOption(pool *radix.Pool) cpn.StrategyOption {
	return poolOption{pool}
}

type poolOption struct {
	pool *radix.Pool
}

func (o poolOption) Apply(p cpn.Strategy) {
	p.(Pool).SetPool(o.pool)
}

func TypeOption(t reflect.Type) cpn.StrategyOption {
	return typeOption{t}
}

type typeOption struct {
	t reflect.Type
}

func (o typeOption) Apply(p cpn.Strategy) {
	p.(*Pop).t = o.t
}

type UnmarshalFunc func(string, interface{}) error

func UnmarshallerOption(f UnmarshalFunc) cpn.StrategyOption {
	return unmarshallerOption{f}
}

type unmarshallerOption struct {
	f UnmarshalFunc
}

func (o unmarshallerOption) Apply(p cpn.Strategy) {
	p.(*Pop).f = o.f
}
