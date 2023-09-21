package config

var C = &Config{}

type Config struct {
	Env    string
	Secret string

	Gateway struct {
		IssuerAddr string
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
