dev:
	find . -type d | justrun -c 'go install github.com/felixge/quantastic/cmd/quantastic && quantastic' -stdin -i .git
