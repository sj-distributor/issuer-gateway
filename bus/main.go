package main

import (
	"cert-gateway/pkg/driver"
)

func main() {
	driver.NewGrpcServiceAndListen(":9527")
}
