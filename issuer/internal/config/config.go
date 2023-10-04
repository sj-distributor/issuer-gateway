package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	Secret    string
	JWTSecret string

	Logger struct {
		Level    string
		Mode     string
		Path     string
		KeepDays int
		MaxSize  int
	}

	Issuer struct {
		rest.RestConf

		CADirURL string

		CheckExpireWithCron struct {
			Type string
			Cron string
		}

		User struct {
			Name string
			Pass string
		}

		Mysql struct {
			User string
			Pass string
			Host string
			Port string
			DB   string
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
