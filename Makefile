
proto/model:
	protoc -I=. -I=model/proto --go_out=plugins=grpc:$(GOPATH)/src $(shell ls ./model/proto/)

proto/crypto:
	protoc -I=. -I=crypto/proto --go_out=plugins=grpc:$(GOPATH)/src $(shell ls ./crypto/proto/)

proto/service:
	protoc -I=. -I=service/proto --go_out=plugins=grpc:$(GOPATH)/src $(shell ls ./service/proto/)

proto: proto/model proto/crypto proto/service

.phony: proto
	