package rtpn

import "context"

type Option interface {
	Apply(p *P)
}

func WithContextOption(ctx context.Context) Option {
	return contextOption{ctx}
}

type contextOption struct {
	ctx context.Context
}

func (o contextOption) Apply(p *P) {
	p.ctx = o.ctx
}

func WithPlaceOption(place Place) Option {
	return storageOption{place}
}

type storageOption struct {
	place Place
}

func (o storageOption) Apply(p *P) {
	p.place = o.place
}
