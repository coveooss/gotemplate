SHELL:=/bin/bash

pre-commit:
	go fmt ./...
	go test ./... -count 1
	./render-doc

codecov:
	bash ./test.sh
	bash <(curl -s https://codecov.io/bash)

# IMPORTANT:
# call `make doc` to generate the doc rendering used to test gotemplate
# Be sure to validate the rendered files before commiting your code
doc:
	./render-doc

# Starts local doc Hugo server (https://gohugo.io/)
# Before calling this you need to install Hugo
doc-serve:
	./render-doc serve
	
# Used to generate the final Hugo static website. Used in CI
doc-generate:
	./render-doc
	hugo --minify --source docs
	git submodule deinit -f docs/themes/book
	git diff -b -w --ignore-blank-lines --exit-code
