serverpath = .
servertag = dev

server/build/docker:
	docker build $(serverpath) -t taask/server:$(servertag)

server/proto/model:
	protoc -I=$(GOPATH)/src -I=. -I=model/proto --go_out=plugins=grpc:$(GOPATH)/src $(shell ls ./model/proto/)

server/proto/service:
	protoc -I=$(GOPATH)/src -I=. -I=service/proto --go_out=plugins=grpc:$(GOPATH)/src $(shell ls ./service/proto/)

server/proto: proto/model proto/service

.phony: proto
	