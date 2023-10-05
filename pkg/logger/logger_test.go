package logger

import (
	"github.com/zeromicro/go-zero/core/logx"
	"testing"
)

func TestInit(t *testing.T) {
	type args struct {
		conf logx.LogConf
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "can init logger",
			args: struct{ conf logx.LogConf }{conf: logx.LogConf{Level: "debug", Mode: "console", Path: "", KeepDays: 0, MaxSize: 0}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Init(tt.args.conf)
		})
	}
}
