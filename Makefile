serverpath = .
servertag = dev
	
build/server/docker: tag/server/dev
	docker build $(serverpath) -t taask/server:$(shell cat ./taask-server/.build/tag)

push/server/docker:	
	docker push taask/server:$(shell cat ./taask-server/.build/tag)

install/server: build/server/docker
	helm template $(serverpath)/ops/chart \
	--set Tag=$(shell cat ./taask-server/.build/tag) --set HomeDir=$(HOME) \
	| linkerd inject --proxy-bind-timeout 30s - \
	| kubectl apply -f - -n taask

install/server/partner:
	helm template $(serverpath)/ops/chart \
	--set Tag=$(shell cat ./taask-server/.build/tag) --set HomeDir=$(HOME) --set Suffix="-partner" --set PartnerHost="taask-server-manager" \
	| linkerd inject --proxy-bind-timeout 30s - \
	| kubectl apply -f - -n taask

logs/server:
	kubectl logs deployment/taask-server taask-server -n taask -f

logs/server/partner:
	kubectl logs deployment/taask-server-partner taask-server -n taask -f

logs/server/search:
	kubectl logs deployment/taask-server taask-server -n taask -f | grep $(search)

uninstall/server:
	kubectl delete service taask-server-internal -n taask
	kubectl delete service taask-server-manager -n taask
	kubectl delete service taask-server-ingress -n taask
	kubectl delete deployment taask-server -n taask

uninstall/server/partner:
	kubectl delete service taask-server-manager-partner -n taask
	kubectl delete deployment taask-server-partner -n taask

tag/server/dev:
	mkdir -p $(serverpath)/.build
	date +%s | openssl sha256 | base64 | head -c 12 > $(serverpath)/.build/tag

install/secrets:
	kubectl create secret generic taask-server-config --from-file=$(HOME)/.taask/server/config/client-auth.yaml --from-file=$(HOME)/.taask/server/config/runner-auth.yaml -n taask

uninstall/secrets:
	kubectl delete secret taask-server-config -n taask

proto/server/model:
	protoc -I=$(GOPATH)/src -I=. -I=model/proto --go_out=plugins=grpc:$(GOPATH)/src $(shell ls ./model/proto/)

proto/server/auth:
	protoc -I=$(GOPATH)/src -I=. -I=auth/proto --go_out=plugins=grpc:$(GOPATH)/src $(shell ls ./auth/proto/)

proto/server/partner:
	protoc -I=$(GOPATH)/src -I=. -I=partner/proto --go_out=plugins=grpc:$(GOPATH)/src $(shell ls ./partner/proto/)

proto/server/service:
	protoc -I=$(GOPATH)/src -I=. -I=service/proto --go_out=plugins=grpc:$(GOPATH)/src $(shell ls ./service/proto/)

proto/server: proto/server/model proto/server/service proto/server/auth proto/server/partner

.phony: proto
	