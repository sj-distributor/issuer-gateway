package acme

import (
	"github.com/go-acme/lego/v4/acme"
	"github.com/go-jose/go-jose/v3/json"
	"github.com/go-playground/assert/v2"
	"sync"
	"testing"
)

func TestReqCertificate(t *testing.T) {

	acmeProvider := &AcmeProvider{
		MemoryCache: &sync.Map{},
	}

	type args struct {
		CADirURL     string
		accountEmail string
		domains      []string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "request certificate must be failed",
			args: struct {
				CADirURL     string
				accountEmail string
				domains      []string
			}{
				CADirURL:     "https://acme-staging-v02.api.letsencrypt.org/directory",
				accountEmail: "nsgzfei@gmail.com", domains: []string{"a.itst.cn"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := acmeProvider.ReqCertificate(tt.args.CADirURL, tt.args.accountEmail, tt.args.domains...)
			if err != nil {
				bytes, err := json.Marshal(err)
				if err != nil {
					t.Errorf("request certificate is err: [%s]", err)
				}

				dict := &map[string]acme.ProblemDetails{}
				err = json.Unmarshal(bytes, dict)
				if err != nil {
					t.Errorf("request certificate is err: [%s]", err)
				}

				for key, _ := range *dict {
					assert.Equal(t, key, "a.itst.cn")
				}
			}

		})
	}
}
