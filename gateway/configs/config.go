package configs

var C = &Config{}

type Config struct {
	Server struct {
		Url    string
		Secret string
	}
	MustHttps bool
	Acme      struct {
		Email string
	}
}
