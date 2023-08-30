package handler

import (
	"cert-gateway/gateway/internal/configs"
	"cert-gateway/gateway/internal/utils"
	"log"
	"net/http"
	"time"
)

func Http() {
	if configs.C.MustHttps {
		go func() {
			mux := http.NewServeMux()

			mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				httpsURL := "https://" + r.Host + r.URL.Path
				http.Redirect(w, r, httpsURL, http.StatusMovedPermanently)
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
}
