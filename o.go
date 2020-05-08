package rtpn

import "context"

type PlaceOption interface {
	Apply(*P)
}

func WithContext(ctx context.Context) PlaceOption {
	return contextOption{ctx}
}

type contextOption struct {
	ctx context.Context
}

func (o contextOption) Apply(p *P) {
	p.ctx = o.ctx
}

func WithPlace(place Place) PlaceOption {
	return placeOption{place}
}

type placeOption struct {
	place Place
}

func (o placeOption) Apply(p *P) {
	p.place = o.place
}

func IsTermination() PlaceOption {
	return terminationOption{}
}

type terminationOption struct{}

func (o terminationOption) Apply(p *P) {
	p.t = true
}

type TransitionOption interface {
	Apply(*T)
}

func WithFunction(fn Transition) TransitionOption {
	return transitionOption{fn}
}

type transitionOption struct {
	fn Transition
}

func (o transitionOption) Apply(t *T) {
	t.fn = o.fn
}
