SHELL:=/bin/bash

.PHONY: pre-commit
pre-commit:
	go fmt ./...
	go test ./... -count 1
	./render-doc

.PHONY: test-all
test-all:
	go test -race -v ./...

# IMPORTANT:
# call `make doc` to generate the doc rendering used to test gotemplate
# Be sure to validate the rendered files before committing your code
.PHONY: doc
doc:
	./render-doc

# Starts local doc Hugo server (https://gohugo.io/)
# Before calling this you need to install Hugo
.PHONY: doc-serve
doc-serve:
	./render-doc
	git submodule update --init
	cd docs && hugo server

# Used to generate the final Hugo static website. Used in CI
.PHONY: doc-generate
doc-generate:
	./render-doc
	git submodule update --init
	hugo --minify --source docs
	git diff -b -w --ignore-blank-lines --exit-code
