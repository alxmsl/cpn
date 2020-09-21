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

func WithKeep(keep bool) POpt {
	return keepOpt{keep}
}

type keepOpt struct {
	keep bool
}

func (o keepOpt) Apply(p *P) {
	p.k = o.keep
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
