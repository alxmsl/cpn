package http

import (
	"fmt"
	"net/http"
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

func (ctx *RequestContext) String() string {
	return fmt.Sprintf("%s %s\n", ctx.r.Method, ctx.r.UserAgent())
}

func (ctx *RequestContext) Wait() {
	<-ctx.done
}
