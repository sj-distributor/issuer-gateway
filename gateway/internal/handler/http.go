package handler

import (
	"cert-gateway/utils"
	"log"
	"net/http"
	"time"
)

func Http() {
	go func() {
		mux := http.NewServeMux()

		mux.HandleFunc("/", ReverseProxyOrRedirect)

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
