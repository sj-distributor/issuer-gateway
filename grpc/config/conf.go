package conf

type Config struct {
	Env    string
	Secret string

	Sync struct {
		Target string
		Grpc   struct {
			Addr string
		}
	}
}
