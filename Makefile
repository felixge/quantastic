dev:
	@find \
		. \
		-type d \
		\! -path '*/.git/*' \
		\! -path '*/static*' \
		! -path '*/templates*' \
		| justrun \
			-c 'go install github.com/felixge/quantastic/cmd/quantastic && quantastic' \
			-stdin

.PHONY: dev
