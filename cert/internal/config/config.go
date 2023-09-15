package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	Secret    string
	Env       string
	JWTSecret string
	Sync      struct {
		Type    string
		Address string
		User    string
		Pass    string
	}
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
