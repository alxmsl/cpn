package http

import (
	"github.com/alxmsl/cpn"
)

type Handler func(ctx *RequestContext)

func Processor(handler Handler) cpn.Transition {
	return func(mm []*cpn.M) *cpn.M {
		m := mm[0]
		handler(m.Value().(*RequestContext))
		return m
	}
}
