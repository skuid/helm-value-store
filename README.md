# helm-value-store

[![Build Status](https://travis-ci.org/skuid/helm-value-store.svg?branch=master)](https://travis-ci.org/skuid/helm-value-store)
[![https://img.shields.io/badge/godoc-reference-5272B4.svg?style=flat-square](https://img.shields.io/badge/godoc-reference-5272B4.svg?style=flat-square)](http://godoc.org/github.com/skuid/helm-value-store/)


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
$ helm value-store list
UniqueId                              Name                  Namespace    Chart                       Version  Labels                                     Values
6fad4903-58ec-446f-bda4-bd39c4ff96aa  alertmanager          default      skuid/alertmanager          0.1.0    map[region:us-west-2 environment:prod]     1.1K
8795d237-adac-4b91-b55b-bb0f1e258a32  exporter              default      skuid/prom-node-exporter    0.1.0    map[region:us-west-2 environment:prod]     279B
22c8f1e8-82fc-4eb0-b1f9-2c8d50b2df3b  prom1                 default      skuid/prometheus            0.1.2    map[region:us-west-2 environment:prod]     1.1K
ad01e6d4-05ec-4f18-ba6a-87cd49e6be25  alertmanager          default      skuid/alertmanager          0.1.0    map[environment:test region:us-west-2]     0
84c28f16-0bc2-4384-9e21-8077e3320aad  exporter              default      skuid/prom-node-exporter    0.1.0    map[environment:test region:us-west-2]     274B
fa718433-d76e-4edd-b263-9c50246c2f80  prom1                 default      skuid/prometheus            0.1.2    map[environment:test region:us-west-2]     0
080f9a8a-10dd-4c2f-8588-8c3c4980553f  alertmanager          default      skuid/alertmanager          0.1.0    map[region:eu-central-1 environment:prod]  1.3K
49582465-85fd-49ce-9778-4bf9d1162a2e  exporter              default      skuid/prom-node-exporter    0.1.0    map[environment:prod region:eu-central-1]  272B
34754bde-3114-43ca-bb23-1d4e16f12f95  prom1                 default      skuid/prometheus            0.1.2    map[environment:prod region:eu-central-1]  0

$ helm value-store list -s environment=test
UniqueId                              Name                  Namespace    Chart                       Version  Labels
ad01e6d4-05ec-4f18-ba6a-87cd49e6be25  alertmanager          default      skuid/alertmanager          0.1.0    map[environment:test region:us-west-2]     0
84c28f16-0bc2-4384-9e21-8077e3320aad  exporter              default      skuid/prom-node-exporter    0.1.0    map[environment:test region:us-west-2]     274B
fa718433-d76e-4edd-b263-9c50246c2f80  prom1                 default      skuid/prometheus            0.1.2    map[environment:test region:us-west-2]     0
```

Install a release, automatically fetching values from the value store:

```
$ helm value-store install --name alertmanager --selector region=us-west-2 --selector environment=prod
Fetched chart skuid/alertmanager to /var/folders/pr/79r611f576jczk_79lfndzgc0000gn/T/370122778/alertmanager-0.1.0.tgz
Installing Release 6fad4903-58ec-446f-bda4-bd39c4ff96aa alertmanager skuid/alertmanager   0.1.0
Successfully installed release alertmanager!
```

Get values out of the value store

```
$ helm value-store get-values --uuid 6fad4903-58ec-446f-bda4-bd39c4ff96aa
aws_region: us-west-2
configMap:
  pagerduty_key: somekey
  name: alertmanager.config
  slack_api_url: https://hooks.slack.com/services/
image:
  repository: prom/alertmanager
  tag: v0.4.2
mounts:
  configPath: /etc/alertmanager
replicaCount: 1
resources:
  limits:
    cpu: 100.0m
    memory: 128Mi
  requests:
    cpu: 50.0m
    memory: 64Mi
service:
  name: alertmanager
  port: 9098
  servicePort: 80
```

Update a release definition:

```
$ helm value-store update --uuid 6fad4903-58ec-446f-bda4-bd39c4ff96aa -f alertmanager-values.yaml
Update release in release store!
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
A helm plugin for working with Helm Release data

Usage:
  helm value-store [command]

Available Commands:
  create      create a release in the release store
  dump        dump the JSON representation of releases
  delete      delete a release in the release store
  get-values  get the values of a release
  install     install or upgrade a release
  list        list the releases
  load        load a json file of releases
  update      update a release in the release store
  version     print the version number

Use "value-store [command] --help" for more information about a command.
```

## License

MIT License (see [LICENSE](/LICENSE))
