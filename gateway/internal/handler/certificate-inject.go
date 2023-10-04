package handler

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/pygzfei/issuer-gateway/gateway/internal/cache"
)

func CertificateInject(info *tls.ClientHelloInfo) (*tls.Certificate, error) {
	if cert, b := cache.GlobalCache.Get(info.ServerName); b {
		return &cert.TlS, nil
	}
	return nil, errors.New(fmt.Sprintf("%s: certificate not found", info.ServerName))
}
