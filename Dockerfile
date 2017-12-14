FROM alpine:latest

MAINTAINER Micah Hausler, <micah.hausler@skuid.com>

ENV HELM_HOME /home/root/.helm
ENV HELM_HOST tiller.kube-system.svc.cluster.local:44134
ENV TILLER_HOST $HELM_HOST

RUN apk -U add ca-certificates

COPY helm-value-store /bin/helm-value-store

ENV AWS_REGION us-west-2

ENTRYPOINT ["/bin/helm-value-store"]
