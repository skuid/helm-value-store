# helm-value-store

A helm plugin for working with helm deployment values.

## Installation

### Prerequisite

You must have go installed

```bash
brew install go
mkdir -p ~/go/{src,bin,pkg}
export GOPATH=~/go
export PATH="$PATH:~/go/bin"

# Append GOPATH to profile
echo 'export GOPATH=~/go' | tee -a ~/.profile
echo 'export PATH="$PATH:~/go/bin"' | tee -a ~/.profile
```

### Install helm-value-store

```bash
go get github.com/skuid/helm-value-store
```

## Usage

```
$ helm-value-store
A tool loading/backing up AWS Dynamo demo data

Usage:
  helm-value-store [command]

Available Commands:
  install     install a release
  list        list the releases
  load        load a json file of releases
  version     Print the version number

Use "helm-value-store [command] --help" for more information about a command.
```

## Development

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

## License

MIT License (see [LICENSE](/LICENSE))
