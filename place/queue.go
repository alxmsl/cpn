package place

import "github.com/alxmsl/rtpn"

type Queue struct {
	q chan *rtpn.M
}

func NewQueue(size int) *Queue {
	return &Queue{make(chan *rtpn.M, size)}
}

func (b Queue) ReadCh() <-chan *rtpn.M {
	return b.q
}

func (b Queue) WriteCh() chan<- *rtpn.M {
	return b.q
}
