package cpn

import "context"

type POpt interface {
	Apply(*P)
}

func WithContext(ctx context.Context) POpt {
	return contextOpt{ctx}
}

type contextOpt struct {
	ctx context.Context
}

func (o contextOpt) Apply(p *P) {
	p.ctx = o.ctx
}

func WithPlace(place Place) POpt {
	return placeOpt{place}
}

type placeOpt struct {
	place Place
}

func (o placeOpt) Apply(p *P) {
	p.place = o.place
}

type PlaceBuilder func(opts ...PlaceOption) Place
type PlaceOption interface {
	Apply(Place)
}

func WithPlaceBuilder(builder PlaceBuilder, opts ...PlaceOption) POpt {
	return placeBuilderOpt{builder, opts}
}

type placeBuilderOpt struct {
	builder PlaceBuilder
	opts    []PlaceOption
}

func (o placeBuilderOpt) Apply(p *P) {
	p.place = o.builder(o.opts...)
}

func IsFinal() POpt {
	return finalOpt{}
}

type finalOpt struct{}

func (o finalOpt) Apply(p *P) {
	p.t = true
}

func IsInitial() POpt {
	return initialOpt{}
}

type initialOpt struct{}

func (o initialOpt) Apply(p *P) {
	p.i = true
}

type TOpt interface {
	Apply(*T)
}

func WithFunction(fn Transition) TOpt {
	return transitionOpt{fn}
}

type transitionOpt struct {
	fn Transition
}

func (o transitionOpt) Apply(t *T) {
	t.fn = o.fn
}
