# helm-value-store

A helm plugin for working with helm deployment values.

When working with multiple Kubernetes clusters, but using the same helm charts,
you might have slightly different values files that need to be stored somewhere.

This project is an attempt to manage working with multiple `values.yaml` files for
nearly identitcal deployments.

The only backing store is currently DynamoDB, but other backends such as etcd, consul,
or Vault could easily be implemented.


## Example

List values from your backend

```
$ helm value-store  list
UniqueId                              Name                  Namespace    Chart                       Version  Labels
8795d237-adac-4b91-b55b-bb0f1e258a32  exporter              default      prom-node-exporter    0.1.0    map[environment:prod region:us-west-2]
22c8f1e8-82fc-4eb0-b1f9-2c8d50b2df3b  prom1                 default      prometheus            0.1.2    map[environment:prod region:us-west-2]
6fad4903-58ec-446f-bda4-bd39c4ff96aa  alertmanager          default      alertmanager          0.1.0    map[environment:prod region:us-west-2]
fa718433-d76e-4edd-b263-9c50246c2f80  prom1                 default      prometheus            0.1.2    map[environment:test region:us-west-2]
84c28f16-0bc2-4384-9e21-8077e3320aad  exporter              default      prom-node-exporter    0.1.0    map[environment:test region:us-west-2]
ad01e6d4-05ec-4f18-ba6a-87cd49e6be25  alertmanager          default      alertmanager          0.1.0    map[environment:test region:us-west-2]

$ helm value-store  list -s environment=test
UniqueId                              Name                  Namespace    Chart                       Version  Labels
fa718433-d76e-4edd-b263-9c50246c2f80  prom1                 default      prometheus            0.1.2    map[environment:test region:us-west-2]
84c28f16-0bc2-4384-9e21-8077e3320aad  exporter              default      prom-node-exporter    0.1.0    map[environment:test region:us-west-2]
ad01e6d4-05ec-4f18-ba6a-87cd49e6be25  alertmanager          default      alertmanager          0.1.0    map[environment:test region:us-west-2]
```

Install multiple releases:

```
Install Releases
$ helm value-store install --selector environment=test --selector region=us-west-2
Installing releases:

helm install --name alertmanager --namespace default --version 0.1.0 alertmanager
helm install --name prom1 --namespace default --version 0.1.2 prometheus
helm install --name exporter --namespace default --version 0.1.0 prom-node-exporter
```


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

### AWS Prerequisite

You must have the ability to create, read, and write to a DynamoDB table.

Set the proper access key environment variables, or use the
`$HOME/.aws/{config/credentials}` and set the appropriate
`AWS_DEFAULT_PROFILE` environment variable.

### Install helm-value-store

```bash
go get github.com/skuid/helm-value-store
```

### Add the plugin to helm

```bash
mkdir -p $HELM_HOME/pluings/value-store
cat <<EOF > $HELM_HOME/plugins/value-store/plugin.yaml
name: "value-store"
version: "0.1.0"
usage: "Store values in DynamoDB"
ignoreFlags: false
useTunnel: false
command: "helm-value-store"
EOF
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
