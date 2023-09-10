#修改模版
temp:
	cp -f template/handler.tpl ~/.goctl/1.5.4/api/handler.tpl

doc:
	rm -rf ./cert/doc/* && goctl api doc --dir="./cert" --o="./cert/doc"

test:
	go test -v ./...

build:
	docker build -f ./cert/Dockerfile -t cert:v1.0.0 . &&  docker build -f ./gateway/Dockerfile -t gateway:v1.0.0 .
proto:
	cd bus && protoc --go_out=. --go-grpc_out=. pubsub.proto && cd ..
