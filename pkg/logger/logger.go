package logger

import (
	"github.com/zeromicro/go-zero/core/logx"
	"os"
)

func Init(level string, serviceName string) {

	logx.MustSetup(logx.LogConf{
		Level:       level,
		ServiceName: serviceName,
		Mode:        "console", // [console,file,volume]
	})

	logx.SetWriter(logx.NewWriter(os.Stdout))
}
