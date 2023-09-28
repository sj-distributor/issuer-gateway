FROM golang:alpine AS builder

ENV GOPROXY=https://goproxy.cn
ENV CGO_ENABLED 0

RUN apk update --no-cache && apk add --no-cache tzdata

USER root
WORKDIR /build

COPY . .
RUN go mod download

COPY ./conf /app/conf

RUN go build -ldflags="-s -w" -o /app/ig ./cmd/main.go


FROM alpine:latest
RUN apk update --no-cache && apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=builder /app/conf /app/conf
COPY --from=builder /app/ig /app/ig

CMD ["/app/ig"]