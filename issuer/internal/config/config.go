package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	Secret    string
	JWTSecret string

	Logger struct {
		Level string
	}

	Issuer struct {
		rest.RestConf

		CADirURL string

		User struct {
			Name string
			Pass string
		}

		Mysql struct {
			Dns string
			Env string
		}
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
