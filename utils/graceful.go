package utils

import (
	"fmt"
	"gopkg.in/tylerb/graceful.v1"
	"net/http"
	"time"
)

func GraceFul(timeout time.Duration, server *http.Server) *graceful.Server {

	srv := &graceful.Server{
		Timeout: timeout,
		Server:  server,
		ShutdownInitiated: func() {
			fmt.Println("ShutdownInitiated executed...")
		},
	}

	return srv
}
