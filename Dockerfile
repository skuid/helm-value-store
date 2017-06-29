# To build:
# $ docker build -t skuid/helm-value-store .
#
# To run:
# $ docker run skuid/helm-value-store

FROM debian:latest

MAINTAINER Micah Hausler, <micah.hausler@skuid.com>

RUN apt-get update && apt-get install -y ca-certificates

COPY helm-value-store /bin/helm-value-store

ENV AWS_REGION us-west-2

ENTRYPOINT ["/bin/helm-value-store"]
