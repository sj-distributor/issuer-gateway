FROM golang:alpine AS builder
USER root
WORKDIR /build

ENV GOPROXY=https://goproxy.cn
ENV CGO_ENABLED 0

COPY . /build
COPY ./cert/internal/configs/config.yaml /app/internal/configs/config.yaml

RUN go mod download && go mod tidy
RUN go build -ldflags="-s -w" -o /etc/cert /build/cert/main.go


FROM alpine:latest

RUN apk update && apk add ca-certificates

WORKDIR /app
COPY --from=builder /etc/cert /app/cert
COPY --from=builder /app/internal/configs/config.yaml /app/internal/configs/config.yaml

EXPOSE 80 443 9527

CMD ["./cert", "-f", "internal/configs/config.yaml"]
