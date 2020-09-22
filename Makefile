.PHONY: bench
bench:
	go test ./... -check.f='!Test' -bench=. -benchmem

.PHONY: fmt
fmt:
	gofmt -l -w `find . -type f -name '*.go' -not -path "./vendor/*"`
	goimports -l -w `find . -type f -name '*.go' -not -path "./vendor/*" -not -path "./place/redis/*" -not -path "./example/redis/*"`

.PHONY: test
test:
	go test ./... -v -race

%:
	@:
