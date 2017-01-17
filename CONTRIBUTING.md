# Development

Always, always, always run `go fmt ./...` before committing!

### Running the tests

```bash
go get golang.org/x/tools/cmd/cover

make test
```

See the html output of the coverage information

```bash
make test-cover
```

### Updating dependencies

```bash
go get -u github.com/kardianos/govendor

govendor add +external
```

### Linting

Perfect linting is not required, but it is helpful for new people coming to the code.

```
go get -u github.com/golang/lint/golint

golint ./
golint ./render
```
