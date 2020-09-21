package place

import (
	"context"

	"github.com/alxmsl/cpn"
)

type Cancel interface {
	SetCancel(context.CancelFunc)
}

func CancelOption(f context.CancelFunc) cpn.PlaceOption {
	return cancelOption{f}
}

type cancelOption struct {
	f context.CancelFunc
}

func (o cancelOption) Apply(p cpn.Place) {
	p.(Cancel).SetCancel(o.f)
}
