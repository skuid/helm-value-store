FROM alpine:latest

MAINTAINER Micah Hausler, <micah.hausler@skuid.com>

RUN apk -U add ca-certificates

COPY helm-value-store /bin/helm-value-store

ENV AWS_REGION us-west-2

ENTRYPOINT ["/bin/helm-value-store"]
