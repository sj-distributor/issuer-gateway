package handler

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/pygzfei/issuer-gateway/gateway/internal/cache"
)

func CertificateInject(info *tls.ClientHelloInfo) (*tls.Certificate, error) {
	fmt.Println()
	fmt.Println(info.ServerName)
	fmt.Println(info.SupportedProtos)
	fmt.Println()
	if cert, b := cache.GlobalCache.Get(info.ServerName); b {
		return &cert.TlS, nil
	}
	return &tls.Certificate{}, errors.New(fmt.Sprintf("%s: certificate not found", info.ServerName))
}
