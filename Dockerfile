FROM golang:1.13.7-buster

WORKDIR /go/src/github.com/kkeisuke/docker-go-gin

RUN curl -fLo /go/bin/air https://git.io/linux_air \
  && chmod +x /go/bin/air

CMD air