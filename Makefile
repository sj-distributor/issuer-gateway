#修改模版
temp:
	cp -f template/handler.tpl ~/.goctl/1.5.4/api/handler.tpl

doc:
	rm -rf ./issuer/docs/* && goctl api doc --dir="./issuer" --o="./issuer/docs"

test:
	go test -v ./...

build:
	go build -ldflags="-s -w" -o ./ig ./cmd/main.go
proto:
	cd grpc && rm -rf ./pb/* && protoc --go_out=. --go-grpc_out=. pubsub.proto && cd ..
