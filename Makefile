pre-commit:
	go fmt ./...
	go test ./...

install:
	glide install
	go install

# IMPORTANT:
# type: make doc
# to generate the doc rendering used to test gotemplate
# Be sure to validate the rendered files before commiting your code
doc:
	./render-doc

coveralls:
	wget https://raw.githubusercontent.com/coveo/terragrunt/master/scripts/coverage.sh
	@sh ./coverage.sh --coveralls
	rm coverage.sh

html-coverage:
	wget https://raw.githubusercontent.com/coveo/terragrunt/master/scripts/coverage.sh
	@sh ./coverage.sh --html
	rm coverage.sh