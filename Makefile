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

lint:
	for pkg in $(GO_PKGS); do golint $$pkg; done

build: fmt
	go build

test: fmt
	go test -race $(GO_PKGS)

test-cover: fmt
	go test -cover $(GO_PKGS)

vendored:
	# Check if any dependencies are missing
	test $$(govendor list +e |wc -l | awk '{print $$1}') -lt 1

completion: build
	./$(REPO) completion > out.sh
	cp out.sh  /usr/local/etc/bash_completion.d/$(REPO)

docker:
	docker run --rm -v $$(pwd):/go/src/github.com/skuid/$(REPO) -w /go/src/github.com/skuid/$(REPO) golang:1.9-alpine sh -c "apk -U add gcc linux-headers musl-dev && go build -v -a -tags netgo -installsuffix netgo -ldflags '-w'"
	docker build -t skuid/$(REPO) .

clean:
	rm ./$(REPO)
