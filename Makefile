PHONY: run

run:
	go run main.go

build:
	CGO_ENABLED=0 go build -o bin/testresults

run_platform:
	$(MAKE) build
	docker build -t test-results-$(PLATFORM) -f dockerfiles/Dockerfile.$(PLATFORM) .
	docker run --rm -ti test-results-$(PLATFORM) $(CMD)

run_alpine:
	$(MAKE) run_platform PLATFORM=alpine CMD=/bin/ash

run_ubuntu1804:
	$(MAKE) run_platform PLATFORM=ubuntu1804 CMD=/bin/bash

run_ubuntu2004:
	$(MAKE) run_platform PLATFORM=ubuntu2004 CMD=/bin/bash

run_arch:
	$(MAKE) run_platform PLATFORM=arch CMD=/bin/bash
