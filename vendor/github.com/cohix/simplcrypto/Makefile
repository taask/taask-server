proto:
	protoc -I=. -I=proto --go_out=plugins=grpc:$(GOPATH)/src $(shell ls ./proto/)

.PHONY: proto