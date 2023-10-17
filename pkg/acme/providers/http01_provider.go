package providers

import (
	"sync"
)

type DomainMapping struct {
	Token   string
	KeyAuth string
}

type Http01ProviderServer struct {
	mappings *sync.Map
}

func NewHttp01ProviderServer(dict *sync.Map) *Http01ProviderServer {
	return &Http01ProviderServer{
		mappings: dict,
	}
}

// Present starts a web server and makes the token available at `ChallengePath(token)` for web requests.
func (s *Http01ProviderServer) Present(domain, token, keyAuth string) error {
	s.mappings.Store(domain, DomainMapping{Token: token, KeyAuth: keyAuth})
	return nil
}

// CleanUp closes the HTTP server and removes the token from `ChallengePath(token)`.
func (s *Http01ProviderServer) CleanUp(domain, token, keyAuth string) error {
	return nil
}
