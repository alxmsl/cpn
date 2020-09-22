package redis

import (
	"github.com/alxmsl/cpn"
	"github.com/mediocregopher/radix/v3"
)

func KeyOption(k string) cpn.PlaceOption {
	return keyOption{k}
}

type keyOption struct {
	k string
}

func (o keyOption) Apply(p cpn.Place) {
	p.(*Push).key = o.k
}

type MarshalFunc func(interface{}) (string, error)

func MarshallerOption(f MarshalFunc) cpn.PlaceOption {
	return marshallerOption{f}
}

type marshallerOption struct {
	f MarshalFunc
}

func (o marshallerOption) Apply(p cpn.Place) {
	p.(*Push).f = o.f
}

func PoolOption(pool *radix.Pool) cpn.PlaceOption {
	return poolOption{pool}
}

type poolOption struct {
	pool *radix.Pool
}

func (o poolOption) Apply(p cpn.Place) {
	p.(*Push).pool = o.pool
}
