REPO=helm-value-store
SHA = $(shell git rev-parse --short HEAD)

.PHONY: setup fmt test vendored clean

all: test build

setup:
	go get golang.org/x/tools/cmd/cover
	go get -u github.com/golang/lint/golint

fmt:
	go fmt ./...

lint:
	for pkg in $$(go list ./...); do golint $$pkg; done

build: fmt
	go build

test: fmt
	go vet  ./...
	go test -cover  ./...
	go test -race ./...

vendored:
	# Check if any dependencies are missing
	test $$(govendor list +e |wc -l | awk '{print $$1}') -lt 1

completion: build
	./$(REPO) completion > out.sh
	cp out.sh  /usr/local/etc/bash_completion.d/$(REPO)

docker:
	docker run --rm -v $$(pwd):/go/src/github.com/skuid/$(REPO) -w /go/src/github.com/skuid/$(REPO) golang:1.9-alpine sh -c "apk -U add gcc linux-headers musl-dev && go build -v -ldflags '-w -X github.com/skuid/helm-value-store/vendor/github.com/skuid/spec/metrics.commit=$(SHA)'"
	docker build -t skuid/$(REPO) .

clean:
	rm ./$(REPO)
