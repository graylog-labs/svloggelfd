FROM ubuntu:14.04

MAINTAINER Graylog Inc. <hello@graylog.com>

RUN apt-get clean
RUN apt-get update
RUN apt-get install -y build-essential curl wget git
RUN wget https://storage.googleapis.com/golang/go1.5.1.linux-amd64.tar.gz
RUN tar -C /usr/local -xzf go1.5.1.linux-amd64.tar.gz
RUN rm go1.5.1.linux-amd64.tar.gz
RUN mkdir -p /var/lib/go
RUN echo 'export GOROOT=/usr/local/go' >> /etc/profile
RUN echo 'export GOPATH=/var/lib/go' >> /etc/profile
RUN echo 'export PATH=$PATH:$GOROOT/bin' >> /etc/profile

# Remove cached packages and metadata to keep images small.
RUN apt-get clean
