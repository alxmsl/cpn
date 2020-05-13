package http

import (
	"context"
	"net/http"

	"github.com/alxmsl/cpn"
	"github.com/alxmsl/cpn/place"
)

type RequestContext struct {
	done chan struct{}

	r *http.Request
	w http.ResponseWriter
}

func (ctx *RequestContext) Flush() {
	close(ctx.done)
}

func (ctx *RequestContext) Response() http.ResponseWriter {
	return ctx.w
}

func (ctx *RequestContext) Request() *http.Request {
	return ctx.r
}

func (ctx *RequestContext) Wait() {
	<-ctx.done
}

type HttpHandler func(ctx *RequestContext)

func HttpProcessor(handler HttpHandler) cpn.Transition {
	return func(mm []*cpn.M) *cpn.M {
		m := mm[0]
		handler(m.Value().(*RequestContext))
		return m
	}
}

type HttpResponse struct {
	cpn.Place
}

func NewHttpResponse() *HttpResponse {
	return &HttpResponse{place.NewBlock()}
}

func (p *HttpResponse) Run() {
	for m := range p.Out() {
		m.Value().(*RequestContext).Flush()
	}
}

type HttpRequest struct {
	cpn.Place

	addr, pattern string
	cancel        context.CancelFunc
}

func NewHttpRequest(addr, pattern string, cancel context.CancelFunc) *HttpRequest {
	return &HttpRequest{
		place.NewBlock(),
		addr, pattern,
		cancel,
	}
}

func (p *HttpRequest) Run() {
	http.HandleFunc(p.pattern, func(w http.ResponseWriter, r *http.Request) {
		ctx := &RequestContext{
			done: make(chan struct{}),
			r:    r,
			w:    w,
		}
		p.In() <- cpn.NewM(ctx)
		ctx.Wait()
	})
	if err := http.ListenAndServe(p.addr, nil); err != http.ErrServerClosed {
		p.cancel()
		return
	}
}
