package handler

import (
	"log"
	"net/http"
	"time"

	"github.com/pygzfei/issuer-gateway/gateway/internal/config"
	"github.com/pygzfei/issuer-gateway/utils"
)

func Http(c *config.Config) {
	go func() {

		mux := http.NewServeMux()

		mux.HandleFunc("/", HttpMiddleware(ReverseProxyOrRedirect))
		mux.HandleFunc("/.well-known/acme-challenge/", HttpMiddleware(AcceptChallenge(c)))

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
