FROM golang:alpine AS builder

ENV GOPROXY=https://goproxy.cn
ENV CGO_ENABLED 0

RUN apk update --no-cache && apk add --no-cache tzdata

USER root
WORKDIR /build


ADD go.mod .
ADD go.sum .
RUN go mod download
COPY . .
COPY ./gateway/etc /app/etc
RUN go build -ldflags="-s -w" -o /app/gateway /build/gateway/gateway.go


FROM scratch

WORKDIR /app
COPY --from=builder /app/gateway /app/gateway
COPY --from=builder /app/etc /app/etc

EXPOSE 80 443

CMD ["./gateway", "-f", "etc/config.yaml"]
