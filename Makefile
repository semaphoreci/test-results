PHONY: run

run:
	go run main.go $(arg)

regen:
	go run main.go compile priv/parsers/generic/in.xml priv/parsers/generic/out.json
	go run main.go compile priv/parsers/rspec/in.xml priv/parsers/rspec/out.json
	go run main.go compile priv/parsers/exunit/in.xml priv/parsers/exunit/out.json
	go run main.go compile priv/parsers/golang/in.xml priv/parsers/golang/out.json
	go run main.go compile priv/parsers/phpunit/in.xml priv/parsers/phpunit/out.json
	go run main.go compile priv/parsers/embedded/in.xml priv/parsers/embedded/out.json
	go run main.go compile priv/merging priv/merging/out.json
	go run main.go compile priv/merging priv/parsers priv/workflow/out.json

test:
	gotestsum ./...

test.watch:
	gotestsum --watch ./...

build:
	go build -o bin/test-results

build.windows:
	CGO_ENABLED=0 GOOS=windows go build -o bin/test-results

release.major:
	git fetch --tags
	latest=$$(git tag | sort --version-sort | tail -n 1); new=$$(echo $$latest | cut -c 2- | awk -F '.' '{ print "v" $$1+1 ".0.0" }');          echo $$new; git tag $$new; git push origin $$new

release.minor:
	git fetch --tags
	latest=$$(git tag | sort --version-sort | tail -n 1); new=$$(echo $$latest | cut -c 2- | awk -F '.' '{ print "v" $$1 "." $$2 + 1 ".0" }');  echo $$new; git tag $$new; git push origin $$new

release.patch:
	git fetch --tags
	latest=$$(git tag | sort --version-sort | tail -n 1); new=$$(echo $$latest | cut -c 2- | awk -F '.' '{ print "v" $$1 "." $$2 "." $$3+1 }'); echo $$new; git tag $$new; git push origin $$new

SECURITY_TOOLBOX_BRANCH ?= master
SECURITY_TOOLBOX_TMP_DIR ?= /tmp/security-toolbox

check.prepare:
	rm -rf $(SECURITY_TOOLBOX_TMP_DIR)
	git clone git@github.com:renderedtext/security-toolbox.git $(SECURITY_TOOLBOX_TMP_DIR) && (cd $(SECURITY_TOOLBOX_TMP_DIR) && git checkout $(SECURITY_TOOLBOX_BRANCH) && cd -)

check.static: check.prepare
	docker run -it -v $$(pwd):/app \
		-v $(SECURITY_TOOLBOX_TMP_DIR):$(SECURITY_TOOLBOX_TMP_DIR) \
		registry.semaphoreci.com/ruby:2.7 \
		bash -c 'cd /app && $(SECURITY_TOOLBOX_TMP_DIR)/code --language go -d'

check.deps: check.prepare
	docker run -it -v $$(pwd):/app \
		-v $(SECURITY_TOOLBOX_TMP_DIR):$(SECURITY_TOOLBOX_TMP_DIR) \
		registry.semaphoreci.com/ruby:2.7 \
		bash -c 'cd /app && $(SECURITY_TOOLBOX_TMP_DIR)/dependencies -d --language go'
