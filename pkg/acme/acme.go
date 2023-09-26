package acme

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"fmt"

	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/challenge/http01"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/registration"
)

// AcmeAccount acme account
type AcmeAccount struct {
	Email        string
	Registration *registration.Resource
	key          crypto.PrivateKey
}

func (u *AcmeAccount) GetEmail() string {
	return u.Email
}
func (u AcmeAccount) GetRegistration() *registration.Resource {
	return u.Registration
}
func (u *AcmeAccount) GetPrivateKey() crypto.PrivateKey {
	return u.key
}

func ReqCertificate(CADirURL, accountEmail string, domains ...string) (*certificate.Resource, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	acmeAccount := AcmeAccount{
		Email: accountEmail,
		key:   privateKey,
	}

	config := lego.NewConfig(&acmeAccount)
	config.CADirURL = CADirURL

	client, err := lego.NewClient(config)

	if err != nil {
		return nil, err
	}

	// 设置http01验证
	err = client.Challenge.SetHTTP01Provider(http01.NewProviderServer("", "10086"))

	if err != nil {
		return nil, err
	}

	//  注册用户
	reg, err := client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
	if err != nil {
		return nil, err
	}
	acmeAccount.Registration = reg

	fmt.Println("-- 开始申请证书 --")

	// 创建证书
	request := certificate.ObtainRequest{
		Domains: []string{domains[0]},
		Bundle:  true,
	}

	certificates, err := client.Certificate.Obtain(request)
	if err != nil {
		return nil, err
	}
	fmt.Println("-- 成功申请 --")

	return certificates, nil
}
