package main

import (
	"context"
	"os"

	"github.com/alxmsl/cpn"
	"github.com/alxmsl/cpn/place"
	"github.com/alxmsl/cpn/place/http"
	"github.com/alxmsl/cpn/place/io"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	n := cpn.NewPN()
	n.P("req",
		cpn.WithContext(ctx),
		cpn.WithPlaceBuilder(http.NewRequest,
			http.AddressOption("127.0.0.1:8080"),
			http.PatternOption("/"),
			place.CancelOption(cancel),
		),
	)
	n.T("echo",
		cpn.WithFunction(http.Processor(func(ctx *http.RequestContext) {
			_ = ctx.Request().Write(ctx.Response())
		})),
	)
	n.P("res",
		cpn.WithContext(context.Background()),
		cpn.WithPlace(http.NewResponse()),
	)
	n.P("log",
		cpn.WithContext(context.Background()),
		cpn.WithPlace(io.NewWriter(io.WriterOption(os.Stdout))),
	)
	n.
		PT("req", "echo").
		TP("echo", "log").
		TP("echo", "res").
		Run()
	<-ctx.Done()
}
