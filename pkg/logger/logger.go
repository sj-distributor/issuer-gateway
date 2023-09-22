package logger

import (
	"github.com/zeromicro/go-zero/core/logx"
	"os"
)

func Init(env string) {

	logx.MustSetup(logx.LogConf{
		Level: env,       // [debug,info,error,severe]
		Mode:  "console", // [console,file,volume]
	})

	logx.SetWriter(logx.NewWriter(os.Stdout))
}
