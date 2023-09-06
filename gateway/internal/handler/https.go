package handler

import (
	"cert-gateway/gateway/internal/cache"
	"crypto/tls"
	"log"
	"net/http"
)

func Https() *http.Server {

	err := cache.GlobalCache.Sync()
	if err != nil {
		log.Fatalln(err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", ReverseProxyHandler)

	server := &http.Server{
		Addr:    ":443",
		Handler: mux,
	}

	server.TLSConfig = &tls.Config{
		GetCertificate: CertificateInject(),
	}

	return server
}
