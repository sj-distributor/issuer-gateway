package handler

import (
	"github.com/pygzfei/issuer-gateway/gateway/internal/config"
	"github.com/pygzfei/issuer-gateway/utils"
	"log"
	"net/http"
	"time"
)

func Http(c *config.Config) {
	go func() {
		mux := http.NewServeMux()

		mux.HandleFunc("/", ReverseProxyOrRedirect)
		mux.HandleFunc("/.well-known/acme-challenge/", AcceptChallenge(c))
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
