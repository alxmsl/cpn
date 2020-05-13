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

func NewPrint(w io.Writer) *Print {
	return &Print{NewBlock(), w}
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
