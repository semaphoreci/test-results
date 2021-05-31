PHONY: run

run:
	go run main.go $(arg)

regen:
	go run main.go compile samples/parsers/generic/in.xml samples/parsers/generic/out.json

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