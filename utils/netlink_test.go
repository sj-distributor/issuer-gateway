package utils

import (
	"fmt"
	"github.com/go-playground/assert/v2"
	"regexp"
	"testing"
)

func TestGetLocalId(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{name: "test GetLocalIP", want: "123"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetLocalIP()
			fmt.Printf("localip: %s \n", got)
			if got == tt.want {
				t.Errorf("GetLocalIP() = %v, want %v", got, tt.want)
			}

			pattern := `^(\d{1,3}\.){3}\d{1,3}$`

			// 编译正则表达式
			regex := regexp.MustCompile(pattern)

			// 使用正则表达式进行匹配
			isValidIP := regex.MatchString(got)

			assert.NotEqual(t, got, tt.want)
			assert.IsEqual(isValidIP, true)
		})
	}
}
