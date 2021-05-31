PHONY: run

run:
	go run main.go $(arg)

regen:
	go run main.go compile priv/parsers/generic/in.xml priv/parsers/generic/out.json
	go run main.go compile priv/parsers/rspec/in.xml priv/parsers/rspec/out.json
	go run main.go compile priv/parsers/exunit/in.xml priv/parsers/exunit/out.json
	go run main.go compile priv/parsers/golang/in.xml priv/parsers/golang/out.json

	go run main.go compile priv/merging priv/merging/out.json

	go run main.go compile priv priv/workflow/out.json

test:
	gotestsum ./...

build:
	CGO_ENABLED=0 go build -o bin/test-results

release.major:
	git fetch --tags
	latest=$$(git tag | sort --version-sort | tail -n 1); new=$$(echo $$latest | cut -c 2- | awk -F '.' '{ print "v" $$1+1 ".0.0" }');          echo $$new; git tag $$new; git push origin $$new

release.minor:
	git fetch --tags
	latest=$$(git tag | sort --version-sort | tail -n 1); new=$$(echo $$latest | cut -c 2- | awk -F '.' '{ print "v" $$1 "." $$2 + 1 ".0" }');  echo $$new; git tag $$new; git push origin $$new

release.patch:
	git fetch --tags
	latest=$$(git tag | sort --version-sort | tail -n 1); new=$$(echo $$latest | cut -c 2- | awk -F '.' '{ print "v" $$1 "." $$2 "." $$3+1 }'); echo $$new; git tag $$new; git push origin $$new