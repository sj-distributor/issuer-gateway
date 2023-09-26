package conf

type Config struct {
	Secret string
	Logger struct {
		Level string
	}

	Sync struct {
		Target     string
		GrpcServer struct {
			Port string
		}
	}
}
