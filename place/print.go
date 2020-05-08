package place

import (
	"fmt"
	"io"
)

type Print struct {
	*Block
	w io.Writer
}

func NewPrint(w io.Writer) *Print {
	return &Print{NewBlock(), w}
}

func (p *Print) Run() *Print {
	go func() {
		for m := range p.Out() {
			if str, ok := m.Value().(fmt.Stringer); ok {
				_, _ = fmt.Fprint(p.w, str)
				continue
			}
			_, _ = fmt.Fprintf(p.w, "%#v\n", m.Value())
		}
	}()
	return p
}
