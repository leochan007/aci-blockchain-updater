FROM golang

LABEL MAINTAINER leo chan <leochan007@163.com>

ENV DEBIAN_FRONTEND noninteractive

RUN go get github.com/mongodb/mongo-go-driver/mongo

COPY core /go/src/github.com/leochan007/aci-blockchain-updater/core

COPY utils /go/src/github.com/leochan007/aci-blockchain-updater/utils

WORKDIR /go/src/github.com/leochan007/aci-blockchain-updater/

RUN go build -o /root/updater core/main.go

RUN rm -rf /go/src/github.com/leochan007/aci-blockchain-updater/

WORKDIR /root
