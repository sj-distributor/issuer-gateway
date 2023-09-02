package request

type CreateCert struct {
	Domain string `json:"domain"`
	Email  string `json:"email"`
}
