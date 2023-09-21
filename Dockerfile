FROM golang:alpine AS builder

ENV GOPROXY=https://goproxy.cn
ENV CGO_ENABLED 0

RUN apk update --no-cache && apk add --no-cache tzdata

USER root
WORKDIR /build

RUN go mod download
COPY . .
COPY ./gateway/etc /app/gateway
COPY ./issuer/etc /app/issuer

RUN go build -ldflags="-s -w" -o ./ig ./cmd/main.go


FROM scratch

WORKDIR /app
COPY --from=builder /app/gateway /app/gateway
COPY --from=builder /app/etc /app/etc

EXPOSE 80 443

CMD ["./ig", "-f", "etc/config.yaml"]
