dev:
	@find \
		. \
		-type f \
		\! -path '*/.*' \
		\! -path '*/static*' \
		! -path '*/templates*' \
		| justrun \
			-c 'go install github.com/felixge/quantastic/cmd/quantastic && quantastic' \
			-stdin

.PHONY: dev
