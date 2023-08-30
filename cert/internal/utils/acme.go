package utils

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/challenge/http01"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/registration"
	"log"
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

func ReqCertificate(accountEmail string, domains ...string) (*certificate.Resource, error) {

	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}

	acmeAccount := AcmeAccount{
		Email: accountEmail,
		key:   privateKey,
	}

	config := lego.NewConfig(&acmeAccount)

	config.CADirURL = lego.LEDirectoryStaging
	config.Certificate.KeyType = certcrypto.RSA2048

	client, err := lego.NewClient(config)
	log.Println("client-----", client, err)

	if err != nil {
		return nil, err
	}

	// 设置http01验证
	err = client.Challenge.SetHTTP01Provider(http01.NewProviderServer("", "80"))
	log.Println("SetHTTP01Provider-----", err)

	if err != nil {
		return nil, err
	}

	//  注册用户
	reg, err := client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
	log.Println("Registration-----", reg, err)
	if err != nil {
		return nil, err
	}
	acmeAccount.Registration = reg

	log.Println(fmt.Printf("%#v\n", reg))

	log.Println(fmt.Println("-"))

	log.Println(fmt.Println("-- 开始申请证书 --"))
	// 创建证书
	request := certificate.ObtainRequest{
		Domains: domains,
		Bundle:  true,
	}
	certificates, err := client.Certificate.Obtain(request)
	if err != nil {
		return nil, err
	}
	log.Println(fmt.Println("-- 开始申请结束 --"))

	log.Println(fmt.Printf("%#v\n", certificates))
	return certificates, nil
}
