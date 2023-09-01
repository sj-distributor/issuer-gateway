package utils

import (
	"cert-gateway/cert/internal/configs"
	"crypto/rand"
	"crypto/rsa"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/lego"
)

func Certificate(accountEmail string, domains ...string) (*certificate.Resource, error) {

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	acmeAccount := AcmeAccount{
		Email: accountEmail,
		key:   privateKey,
	}

	config := lego.NewConfig(&acmeAccount)

	config.CADirURL = lego.LEDirectoryProduction
	if configs.C.Acme.Env == "debug" {
		config.CADirURL = lego.LEDirectoryStaging
	}

	//httpClient := &http.Client{}
	//httpClient.Transport = &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	//config.HTTPClient = &http.Client{}

	return nil, nil

}

type Http01Provider struct {
}

func (p *Http01Provider) Present(domain, token, keyAuth string) error {

	return nil
}

func (p *Http01Provider) CleanUp(domain, token, keyAuth string) error {

	return nil
}
