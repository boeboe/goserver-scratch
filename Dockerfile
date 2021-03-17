FROM golang:alpine AS builder


WORKDIR /tmp/goserver

ADD ./goserver /tmp/goserver

RUN apk update && apk add --no-cache go upx && \
    cd /tmp/goserver && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags '-w -extldflags "-static"' -o /go/bin/goserver main.go && \
    upx /go/bin/goserver

FROM scratch

ARG VERSION=latest

LABEL component="goserver-scratch"
LABEL description="Golang binary goserver in a scratch container"
LABEL version=${VERSION}
LABEL maintainer="Bart Van Bos <bartvanbos@gmail.com>"
LABEL source-repo="https://github.com/boeboe/goserver-scratch"

WORKDIR /
COPY --from=builder /go/bin/goserver /goserver

ENTRYPOINT [ "/goserver" ]

EXPOSE 8080
