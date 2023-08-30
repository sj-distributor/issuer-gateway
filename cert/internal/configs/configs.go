package configs

var C = &Config{}

type Config struct {
	Mysql struct {
		Dns       string
		Migration struct {
			Path     string
			Database string
		}
	}
	Env  string
	Acme struct {
		Email string
	}
}
