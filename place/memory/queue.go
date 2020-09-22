package memory

import "github.com/alxmsl/cpn"

const defaultLength = 1

func LengthOption(length int) cpn.PlaceOption {
	return lengthOption{length}
}

type lengthOption struct {
	l int
}

func (o lengthOption) Apply(p cpn.Place) {
	p.(*Queue).l = o.l
}

type Queue struct {
	q chan *cpn.M
	l int
}

func NewQueue(opts ...cpn.PlaceOption) cpn.Place {
	p := &Queue{
		l: defaultLength,
	}
	for _, o := range opts {
		o.Apply(p)
	}
	p.q = make(chan *cpn.M, p.l)
	return p
}

func (b Queue) Run() {}

func (q Queue) In() chan<- *cpn.M {
	return q.q
}

func (q Queue) Out() <-chan *cpn.M {
	return q.q
}
