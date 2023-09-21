package conf

type Config struct {
	Env    string
	Secret string

	Sync struct {
		Target     string
		GrpcServer struct {
			Port string
		}
	}
}
