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
# Before calling this you need to:
# - Install Hugo
# - call `git submodule update --init` to fetch the Hugo theme
doc-serve:
	./render-doc
	cd docs && hugo server

# Used to generate the final Hugo static website. Used in CI
doc-generate:
	./render-doc && \
	git submodule update --init && \
	hugo --minify --source docs && \
	git diff -b -w --ignore-blank-lines --exit-code