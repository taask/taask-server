
proto/model:
	protoc -I=src -I=src/model/proto --go_out=plugins=grpc:src/model $(shell ls ./src/model/proto/)

proto/crypto:
	protoc -I=src -I=src/crypto/proto --go_out=plugins=grpc:src/crypto $(shell ls ./src/crypto/proto/)

proto/service:
	protoc -I=src -I=src/service/proto --go_out=plugins=grpc:src/service $(shell ls ./src/service/proto/)

proto: proto/model proto/crypto proto/service

.phony: proto
	