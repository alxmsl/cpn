package http

import "github.com/alxmsl/cpn"

type Response struct {
	chin  chan *cpn.M
	chout chan *cpn.M
}

func NewResponse() *Response {
	return &Response{
		chin:  make(chan *cpn.M),
		chout: make(chan *cpn.M),
	}
}

func (p *Response) In() chan<- *cpn.M {
	return p.chin
}

func (p *Response) Out() <-chan *cpn.M {
	return p.chout
}

func (p *Response) Run() {
	defer close(p.chout)
	for m := range p.chin {
		m.Value(m.Word()[len(m.Word())-1], len(m.Word())-1).(*RequestContext).Flush()
	}
}
