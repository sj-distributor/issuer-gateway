package handler

import (
	"net/http"
	"testing"
)

func TestHttpMiddleware(t *testing.T) {

	tests := []struct {
		name string
	}{
		{
			name: "can create instance",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			middleware := HttpMiddleware(func(w http.ResponseWriter, r *http.Request) {

			})
			if middleware == nil {
				t.Errorf("middleware is nil")
			}
		})
	}
}
