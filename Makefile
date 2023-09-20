#修改模版
temp:
	cp -f template/handler.tpl ~/.goctl/1.5.4/api/handler.tpl

doc:
	rm -rf ./issuer/doc/* && goctl api doc --dir="./issuer" --o="./issuer/doc"

test:
	go test -v ./...

build:
	docker build -f ./issuer/Dockerfile -t cert:v1.0.0 . &&  docker build -f ./gateway/Dockerfile -t gateway:v1.0.0 .
proto:
	cd bus && rm -rf ./pb/* && protoc --go_out=. --go-grpc_out=. pubsub.proto && cd ..
