package main

import (
	"context"

	"github.com/alxmsl/cpn"
	"github.com/alxmsl/cpn/place/http"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	n := cpn.NewPN()
	n.P("req",
		cpn.WithContext(ctx),
		cpn.WithPlace(http.NewHttpRequest("127.0.0.1:8080", "/", cancel).Run()),
	)
	n.T("echo",
		cpn.WithFunction(http.HttpProcessor(func(ctx *http.RequestContext) {
			_ = ctx.Request().Write(ctx.Response())
		})),
	)
	n.P("res",
		cpn.WithContext(ctx),
		cpn.WithPlace(http.NewHttpResponse().Run()),
		cpn.IsTermination(),
	)

	n.PT("req", "echo").TP("echo", "res").Run()
	<-ctx.Done()
}
