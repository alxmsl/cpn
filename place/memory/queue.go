package memory

import "github.com/alxmsl/cpn"

type Queue struct {
	q chan *cpn.M
}

func NewQueue(size int) *Queue {
	return &Queue{make(chan *cpn.M, size)}
}

func (b Queue) Run() {}

func (q Queue) In() chan<- *cpn.M {
	return q.q
}

func (q Queue) Out() <-chan *cpn.M {
	return q.q
}
