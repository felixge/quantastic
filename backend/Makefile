name = quantastic
release = 0.0.
commit = $(shell git describe --tags --dirty --always)
gopath = $(shell go env GOPATH)
distdir = dist

dist: client server

client server:
	go install \
		-ldflags "-X main.buildRelease 0.0.0 -X main.buildCommit $(commit)" \
		github.com/felixge/$(name)/backend/cmd/$(name)-$@
	mkdir -p "$(distdir)"
	mv "$(gopath)/bin/$(name)-$@" "$(distdir)"


.PHONY: dist client server
