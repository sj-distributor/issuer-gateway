package config

import "github.com/zeromicro/go-zero/rest"

var C = &Config{}

type Config struct {
	Secret string

	Logger struct {
		Level string
	}

	Gateway struct {
		rest.RestConf
		IssuerService string
	}

	Sync struct {
		Target     string
		GrpcClient struct {
			Listen string
		}
		Redis struct {
			Addrs      []string
			User       string
			Pass       string
			MasterName string
			Db         int
		}
	}
}
