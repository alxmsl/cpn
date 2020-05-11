package memory

import "github.com/alxmsl/cpn"

type Queue struct {
	q chan *cpn.M
}

func NewQueue(size int) *Queue {
	return &Queue{make(chan *cpn.M, size)}
}

func (b Queue) In() chan<- *cpn.M {
	return b.q
}

func (b Queue) Out() <-chan *cpn.M {
	return b.q
}
