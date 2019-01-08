FROM golang:1.11.2 as builder

RUN mkdir -p /go/src/github.com/taask/taask-server
WORKDIR /go/src/github.com/taask/taask-server

COPY . .

RUN go build

FROM debian:stable-slim

RUN mkdir /taask
RUN mkdir -p $HOME/.taask/server/config

COPY --from=builder /go/src/github.com/taask/taask-server/taask-server /taask/taask-server