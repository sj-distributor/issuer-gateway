package config

import "github.com/zeromicro/go-zero/rest"

var C = &Config{}

type Config struct {
	Env    string
	Secret string

	Gateway struct {
		rest.RestConf
		IssuerAddr string
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
