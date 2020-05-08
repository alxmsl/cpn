package main

import (
	"context"

	"github.com/alxmsl/rtpn"
	"github.com/alxmsl/rtpn/place/http"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	n := rtpn.NewPN()
	n.P("req",
		rtpn.WithContext(ctx),
		rtpn.WithPlace(http.NewHttpRequest("127.0.0.1:8080", "/", cancel).Run()),
	)
	n.T("echo",
		rtpn.WithFunction(http.HttpProcessor(func(ctx *http.RequestContext) {
			_ = ctx.Request().Write(ctx.Response())
		})),
	)
	n.P("res",
		rtpn.WithContext(ctx),
		rtpn.WithPlace(http.NewHttpResponse().Run()),
		rtpn.IsTermination(),
	)

	n.PT("req", "echo").TP("echo", "res").Run()
	<-ctx.Done()
}
