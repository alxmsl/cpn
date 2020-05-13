package place

import (
	"fmt"
	"io"

	"github.com/alxmsl/cpn"
)

type Print struct {
	cpn.Place

	w io.Writer
}

func WriterOption(w io.Writer) cpn.PlaceOption {
	return writerOption{w}
}

type writerOption struct {
	w io.Writer
}

func (o writerOption) Apply(p cpn.Place) {
	p.(*Print).w = o.w
}

func NewPrint(opts ...cpn.PlaceOption) cpn.Place {
	p := &Print{Place: NewBlock(opts...)}
	for _, o := range opts {
		o.Apply(p)
	}
	return p
}

func (p *Print) Run() {
	for m := range p.Out() {
		if str, ok := m.Value().(fmt.Stringer); ok {
			_, _ = fmt.Fprint(p.w, str)
			continue
		}
		_, _ = fmt.Fprintf(p.w, "%#v\n", m.Value())
	}
}
