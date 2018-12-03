
proto/model:
	protoc -I=$(GOPATH)/src -I=. -I=model/proto --go_out=plugins=grpc:$(GOPATH)/src $(shell ls ./model/proto/)

proto/service:
	protoc -I=$(GOPATH)/src -I=. -I=service/proto --go_out=plugins=grpc:$(GOPATH)/src $(shell ls ./service/proto/)

proto: proto/model proto/service

.phony: proto
	