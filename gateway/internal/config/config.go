package config

var C = &Config{}

type Config struct {
	Server struct {
		Url    string
		Secret string
	}
}
