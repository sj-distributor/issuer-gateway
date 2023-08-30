package handler

import (
	"cert-gateway/gateway/internal/cache"
	"crypto/tls"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

func Https() *http.Server {

	// 读取证书文件和私钥文件内容
	certPEM, err := os.ReadFile("internal/cert/root.crt")
	if err != nil {
		panic(err)
	}
	keyPEM, err := os.ReadFile("internal/cert/root.key")
	if err != nil {
		panic(err)
	}

	cache.GlobalCache.Set("test.anson.com", &cache.Cert{
		Key:         keyPEM,
		Certificate: certPEM,
	})

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		target, _ := url.Parse("http://192.167.167.167:9527")
		r.Host = target.Host

		proxy := httputil.NewSingleHostReverseProxy(target)
		proxy.ServeHTTP(w, r)
	})

	server := &http.Server{
		Addr:    ":443",
		Handler: mux,
	}

	server.TLSConfig = &tls.Config{
		GetCertificate: func(info *tls.ClientHelloInfo) (*tls.Certificate, error) {
			name := info.ServerName
			if domain, b := cache.GlobalCache.Get(name); b {

				cert := domain.(*cache.Cert)

				pair, err := tls.X509KeyPair(cert.Certificate, cert.Key)
				if err != nil {
					log.Fatalln(err)
				}

				return &pair, nil
			}

			return &tls.Certificate{}, nil
		},
	}

	return server
}
