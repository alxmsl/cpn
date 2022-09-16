package cpn

import "context"

// PlaceOption is an abstraction to define place options
type PlaceOption interface {
	Apply(*P)
}

// WithContext creates an option to use context
func WithContext(ctx context.Context) PlaceOption {
	return contextOpt{ctx}
}

type contextOpt struct {
	ctx context.Context
}

func (o contextOpt) Apply(p *P) {
	p.ctx = o.ctx
}

// WithKeep creates an option to keep token in a terminal place
func WithKeep(keep bool) PlaceOption {
	return keepOpt{keep}
}

type keepOpt struct {
	keep bool
}

func (o keepOpt) Apply(p *P) {
	p.k = o.keep
}

// WithStrategy creates an option to use specific strategy
func WithStrategy(s Strategy) PlaceOption {
	return strategyOpt{s}
}

type strategyOpt struct {
	strategy Strategy
}

func (o strategyOpt) Apply(p *P) {
	p.strategy = o.strategy
}

// StrategyBuilder creates new strategy with required options
type StrategyBuilder func(opts ...StrategyOption) Strategy

// StrategyOption is an abstraction to define strategy option
type StrategyOption interface {
	Apply(Strategy)
}

// WithStrategyBuilder returns an option to create strategy with required options
func WithStrategyBuilder(builder StrategyBuilder, opts ...StrategyOption) PlaceOption {
	return strategyBuilderOpt{builder, opts}
}

type strategyBuilderOpt struct {
	builder StrategyBuilder
	opts    []StrategyOption
}

func (o strategyBuilderOpt) Apply(p *P) {
	p.strategy = o.builder(o.opts...)
}

// TransitionOption is an abstraction to define a transition option
type TransitionOption interface {
	Apply(*T)
}

// WithTransformation return a transition option to use specified transformation
func WithTransformation(fn Transformation) TransitionOption {
	return transformationOpt{fn}
}

type transformationOpt struct {
	transformation Transformation
}

func (o transformationOpt) Apply(t *T) {
	t.transformation = o.transformation
}
