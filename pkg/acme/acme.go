package acme

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/challenge/http01"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/registration"
	"github.com/zeromicro/go-zero/core/logx"
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

func ReqCertificate(env, accountEmail string, domains ...string) (*certificate.Resource, error) {
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
	if env == "debug" {
		config.CADirURL = lego.LEDirectoryStaging
	}

	client, err := lego.NewClient(config)
	logx.Info("client-----", client, err)

	if err != nil {
		return nil, err
	}

	// 设置http01验证
	err = client.Challenge.SetHTTP01Provider(http01.NewProviderServer("", "5001"))
	logx.Info("SetHTTP01Provider-----", err)

	if err != nil {
		return nil, err
	}

	//  注册用户
	reg, err := client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
	logx.Info("Registration-----", reg, err)
	if err != nil {
		return nil, err
	}
	acmeAccount.Registration = reg

	logx.Info("-- 开始申请证书 --")

	// 创建证书
	request := certificate.ObtainRequest{
		Domains: domains,
		Bundle:  true,
	}

	certificates, err := client.Certificate.Obtain(request)
	if err != nil {
		return nil, err
	}
	logx.Info("-- 开始申请结束 --")

	logx.Infof("%#v", certificates)

	return certificates, nil
}
