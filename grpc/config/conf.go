package conf

type Config struct {
	Secret string
	Logger struct {
		Level    string
		Mode     string
		Path     string
		KeepDays int
		MaxSize  int
	}

	Sync struct {
		Target     string
		GrpcServer struct {
			Port string
		}
	}
}
