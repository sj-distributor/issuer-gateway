package cache

type Cert struct {
	Certificate []byte `json:"certificate"`
	Key         []byte `json:"key"`
}
