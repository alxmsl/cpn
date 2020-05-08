package memory

import "github.com/alxmsl/rtpn"

type Queue struct {
	q chan *rtpn.M
}

func NewQueue(size int) *Queue {
	return &Queue{make(chan *rtpn.M, size)}
}

func (b Queue) In() chan<- *rtpn.M {
	return b.q
}

func (b Queue) Out() <-chan *rtpn.M {
	return b.q
}
