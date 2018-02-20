pre-commit:
	go fmt ./...
	go test ./...

install:
	go install

# IMPORTANT:
# type: make doc
# to generate the doc rendering used to test gotemplate
# Be sure to validate the rendered files before commiting your code
doc:
	./render-doc
	