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
COPY ./cert/etc /app/etc
COPY ./cert/internal/database/migration /app/internal/database/migration

RUN go build -ldflags="-s -w" -o /app/issuer /build/issuer/issuer.go


FROM scratch

WORKDIR /app
COPY --from=builder /app/issuer /app/issuer
COPY --from=builder /app/etc /app/etc
COPY --from=builder /app/internal/database/migration /app/internal/database/migration

EXPOSE 80 5001 9527

CMD ["./issuer", "-f", "etc/issuer-api.yaml"]
