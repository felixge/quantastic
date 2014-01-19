release = 0.0.
commit = $(shell git describe --tags --dirty --always)
gopath = $(shell go env GOPATH)
distdir = dist

dist:
	go install \
		-ldflags "-X main.buildRelease 0.0.0 -X main.buildCommit $(commit)" \
		github.com/felixge/quantastic/backend/cmd/quantastic
	mkdir -p "$(distdir)"
	mv "$(gopath)/bin/quantastic" "$(distdir)"

.PHONY: dist
