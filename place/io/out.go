package io

import (
	"fmt"
	"io"

	"github.com/alxmsl/cpn"
)

func WriterOption(w io.Writer) cpn.PlaceOption {
	return writerOption{w}
}

type writerOption struct {
	w io.Writer
}

func (o writerOption) Apply(p cpn.Place) {
	p.(*Writer).w = o.w
}

type Writer struct {
	chin  chan *cpn.M
	chout chan *cpn.M

	w io.Writer
}

func (p *Writer) In() chan<- *cpn.M {
	return p.chin
}

func (p *Writer) Out() <-chan *cpn.M {
	return p.chout
}

func NewWriter(opts ...cpn.PlaceOption) cpn.Place {
	p := &Writer{
		chin:  make(chan *cpn.M),
		chout: make(chan *cpn.M),
	}
	for _, o := range opts {
		o.Apply(p)
	}
	return p
}

func (p *Writer) Run() {
	defer close(p.chout)
	for m := range p.chin {
		if str, ok := m.Value().(fmt.Stringer); ok {
			_, _ = fmt.Fprint(p.w, str)
		} else {
			_, _ = fmt.Fprintf(p.w, "%#v\n", m.Value())
		}
		p.chout <- m
	}
}
