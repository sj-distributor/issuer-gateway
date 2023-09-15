package config

var C = &Config{}

type Config struct {
	IssuerUrl string
	Secret    string
	Sync      struct {
		Type    string
		Address string
		User    string
		Pass    string
	}
}
