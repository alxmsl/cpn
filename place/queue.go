package place

import "github.com/alxmsl/pn"

type Queue struct {
	q chan *pn.M
}

func NewQueue(size int) *Queue {
	return &Queue{make(chan *pn.M, size)}
}

func (b Queue) ReadCh() <-chan *pn.M {
	return b.q
}

func (b Queue) WriteCh() chan<- *pn.M {
	return b.q
}
