DEPS = $(wildcard */*.go)
VERSION = $(shell git describe --always --dirty)
COMMIT_SHA1 = $(shell git rev-parse HEAD)
BUILD_DATE = $(shell date +%Y-%m-%d)
OS = linux
ARCH = amd64

all: vendor lint vet prometheus-puppetdb-exporter

prometheus-puppetdb-exporter: main.go $(DEPS)
	GO111MODULE=on CGO_ENABLED=0 GOOS=linux \
	  go build -a \
		  -ldflags="-X main.version=$(VERSION) -X main.commitSha1=$(COMMIT_SHA1) -X main.buildDate=$(BUILD_DATE)" \
	    -installsuffix cgo -o $@ $<
	strip $@

release: prometheus-puppetdb-exporter-$(VERSION).$(OS)-$(ARCH).tar.gz

%.tar.gz: prometheus-puppetdb-exporter LICENSE
	tar cvzf $@ --transform 's,^,$*/,' $^

clean:
	rm -f prometheus-puppetdb-exporter

lint:
	@ go get -v golang.org/x/lint/golint
	@for file in $$(git ls-files '*.go' | grep -v '_workspace/' | grep -v 'vendor/'); do \
		export output="$$(golint $${file} | grep -v 'type name will be used as docker.DockerInfo')"; \
		[ -n "$${output}" ] && echo "$${output}" && export status=1; \
	done; \
	exit $${status:-0}

vet: main.go
	go vet $<

vendor:
	go mod vendor

.PHONY: all lint vet clean vendor
