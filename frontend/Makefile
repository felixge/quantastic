dev:
	@find \
		. \
		-type d \
		\! -path '*/.*' \
		\! -path './dist' \
		| justrun \
			-c 'rm -f dist/index.html; make dist' \
			-stdin

dist: dist/styles.css dist/app.js dist/index.html

dist/index.html: index.html
	cp $^ $@

dist/styles.css: css/main.css
	cp css/main.css dist/styles.css

dist/app.js: **/*.ts
	tsc --outDir dist/js -m commonjs app/main.ts
	browserify dist/js/main.js -o $@
	rm -rf dist/js

.PHONY: dev dist
