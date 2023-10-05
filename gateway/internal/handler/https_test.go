package handler

import (
	"context"
	"github.com/pygzfei/issuer-gateway/utils"
	"gopkg.in/tylerb/graceful.v1"
	"testing"
	"time"
)

func TestHttps(t *testing.T) {

	tests := []struct {
		name string
		want *graceful.Server
	}{
		{
			name: "https server can startup",
			want: utils.GraceFul(time.Duration(1)*time.Second, Https()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			go tt.want.ListenAndServe()

			time.Sleep(2 * time.Second)

			_ = tt.want.Shutdown(context.Background())
			
		})
	}
}
