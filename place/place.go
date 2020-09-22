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

type Erring interface {
	SetErrs(chan<- error)
}

func ErrorsOutOption(errs chan<- error) cpn.PlaceOption {
	return errorsOutOption{errs}
}

type errorsOutOption struct {
	errs chan<- error
}

func (o errorsOutOption) Apply(p cpn.Place) {
	p.(Erring).SetErrs(o.errs)
}

var ErrorProcess = func(errors <-chan error, f func(error)) {
	for {
		select {
		case err := <-errors:
			f(err)
		}
	}
}
