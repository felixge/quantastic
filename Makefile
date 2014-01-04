dev:
	@find \
		. \
		-type f \
		\! -path '*/.*' \
		\! -path '*/static*' \
		\! -path '*/templates*' \
		| justrun \
			-c 'go build -o quantastic github.com/felixge/quantastic/cmd/server && ./quantastic' \
			-stdin

.PHONY: dev
