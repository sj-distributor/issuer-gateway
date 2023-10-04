package logger

import (
	"github.com/zeromicro/go-zero/core/logx"
	"os"
)

func Init(conf logx.LogConf) {

	logx.MustSetup(conf)

	logx.SetWriter(logx.NewWriter(os.Stdout))
}
