package utils

import (
	"fmt"
	"github.com/go-playground/assert/v2"
	"testing"
)

func TestGetLocalId(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{name: "test GetLocalId", want: "123"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetLocalIP()
			fmt.Printf("localip: %s \n", got)
			if got == tt.want {
				t.Errorf("GetLocalId() = %v, want %v", got, tt.want)
			}
			assert.NotEqual(t, got, tt.want)
		})
	}
}
