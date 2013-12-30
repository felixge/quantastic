dev:
	git ls-files | grep -v -E 'css|html' | justrun -c 'go install github.com/felixge/quantastic/cmd/quantastic && quantastic' -stdin
