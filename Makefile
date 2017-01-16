REPO=helm-value-store
SHA = $(shell git rev-parse --short HEAD)
GO_PKGS=$$(go list ./... | grep -v vendor)

.PHONY: setup fmt test test-cover vendored clean

all: test build

setup:
	go get golang.org/x/tools/cmd/cover
	go get -u github.com/kardianos/govendor
	go get -u github.com/golang/lint/golint

fmt:
	go fmt $(GO_PKGS)

build: fmt
	go build

test: fmt
	go test -race $(GO_PKGS)

test-cover: fmt
	go test -cover $(GO_PKGS)

vendored:
	# Check if any dependencies are missing
	test $$(govendor list +e |wc -l | awk '{print $$1}') -lt 1

clean:
	rm ./$(REPO)
