package utils

import (
	"github.com/go-playground/assert/v2"
	"testing"
)

func TestId(t *testing.T) {
	tests := []struct {
		name string
		want int64
	}{
		{name: "can generate id", want: 123},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Id()
			if got == tt.want {
				t.Errorf("Id() = %v, want %v", got, tt.want)
			}

			assert.NotEqual(t, got, tt.want)
		})
	}
}
