package cache

import "crypto/tls"

type Cert struct {
	Id          uint64 `json:"id"`
	Domain      string `json:"domain"`
	Certificate string `json:"certificate"`
	PrivateKey  string `json:"private_key"`
	Target      string `json:"target"`
	TlS         tls.Certificate
}

type Data struct {
	Certs []Cert `json:"certs"`
}
