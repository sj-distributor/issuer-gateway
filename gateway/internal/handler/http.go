package handler

import (
	"cert-gateway/gateway/internal/cache"
	"cert-gateway/utils"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

func Http() {
	go func() {
		mux := http.NewServeMux()

		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			if cert, b := cache.GlobalCache.Get(r.Host); b {
				if cert.Certificate != "" {
					httpsURL := "https://" + r.Host + r.URL.Path
					http.Redirect(w, r, httpsURL, http.StatusMovedPermanently)
				} else {

					target, err := url.Parse(cert.Target)
					if err != nil {
						w.WriteHeader(http.StatusNotFound)
						return
					}

					r.Host = target.Host

					proxy := httputil.NewSingleHostReverseProxy(target)
					proxy.ServeHTTP(w, r)
				}
				return
			}

			w.WriteHeader(http.StatusNotFound)
		})

		server := &http.Server{
			Addr:    ":80",
			Handler: mux,
		}

		err := utils.GraceFul(time.Minute, server).ListenAndServe()

		if err != nil {
			log.Fatalln(err)
		}
	}()
}
