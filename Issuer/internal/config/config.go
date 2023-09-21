package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	Secret    string
	Env       string
	JWTSecret string

	Issuer struct {
		rest.RestConf

		User struct {
			Name string
			Pass string
		}

		Mysql struct {
			Dns       string
			Migration struct {
				Path string
				Db   string
			}
		}
	}

	Sync struct {
		Target string
		Grpc   struct {
			Addr string
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
