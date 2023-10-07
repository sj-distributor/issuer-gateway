package logger

import (
	"github.com/zeromicro/go-zero/core/logx"
)

func Init(conf logx.LogConf) {
	logx.MustSetup(conf)

	// todo:  the log written to ...
}
