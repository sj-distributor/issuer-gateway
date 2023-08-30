FROM golang:alpine AS builder
USER root
WORKDIR /build

ENV GOPROXY=https://goproxy.cn
ENV CGO_ENABLED 0

COPY . /build
COPY ./cert/internal/configs/ /app/internal/configs/

RUN go mod download && go mod tidy
RUN go build -ldflags="-s -w" -o /etc/cert /build/cert/main.go


FROM scratch

WORKDIR /app
COPY --from=builder /etc/cert /app/cert
COPY --from=builder /app/internal/configs/ /app/internal/configs/

EXPOSE 9527

CMD ["./cert", "-f", "internal/configs/config.yaml"]
