pre-commit:
	go fmt ./...
	go test ./... -count 1
	./render-doc

install:
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

# Starts a jekyll server identical to the one on github pages
# Need ruby and this gem:
# gem install bundler
run-doc-server:
	cd docs && bundle install --path vendor/bundle
	cd docs && bundle exec jekyll serve
