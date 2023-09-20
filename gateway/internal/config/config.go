package config

var C = &Config{}

type Config struct {
	IssuerAddr string
	Secret     string

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
