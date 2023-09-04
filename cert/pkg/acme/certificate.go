package acme

import (
	"cert-gateway/cert/internal/database/entity"
	"cert-gateway/utils"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/jinzhu/copier"
	"time"
)

// GetCertificateExpireTime 获取证书过期时间
func GetCertificateExpireTime(certPEM string) (expire time.Time, err error) {

	// 解码 PEM 格式的证书
	block, _ := pem.Decode([]byte(certPEM))
	if block == nil {
		return time.Now().Add(-24 * time.Hour), errors.New("error decoding certificate PEM")
	}

	// 解析证书
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return time.Now().Add(-24 * time.Hour), errors.New(fmt.Sprintf("Error parsing certificate: %s", err))
	}

	// 获取证书的过期日期 (并且提前3天作为过期时间)
	return cert.NotAfter.Add(-72 * time.Hour), nil
}

func EncryptCertificate(certInfo *certificate.Resource, cert *entity.Cert, secret string) error {

	certificateEncrypt, err := utils.Encrypt(certInfo.Certificate, []byte(secret))

	if err != nil {
		return err
	}
	cert.Certificate = string(certificateEncrypt)

	privateKey, err := utils.Encrypt(certInfo.PrivateKey, []byte(secret))
	if err != nil {
		return err
	}
	cert.PrivateKey = string(privateKey)

	issuerCertificate, err := utils.Encrypt(certInfo.IssuerCertificate, []byte(secret))
	if err != nil {
		return err
	}
	cert.IssuerCertificate = string(issuerCertificate)

	expire, err := GetCertificateExpireTime(string(certInfo.Certificate))
	if err != nil {
		return err
	}
	cert.Expire = expire

	return nil
}

func DecryptCertificate(inputCert *entity.Cert, secret string) (*entity.Cert, error) {

	outputCert := &entity.Cert{}

	err := copier.Copy(outputCert, inputCert)
	if err != nil {
		return nil, err
	}

	certificateDecrypt, err := utils.Decrypt([]byte(outputCert.Certificate), []byte(secret))
	if err != nil {
		return nil, err
	}
	outputCert.Certificate = string(certificateDecrypt)

	privateKey, err := utils.Decrypt([]byte(outputCert.PrivateKey), []byte(secret))
	if err != nil {
		return nil, err
	}
	outputCert.PrivateKey = string(privateKey)

	issuerCertificate, err := utils.Decrypt([]byte(outputCert.IssuerCertificate), []byte(secret))
	if err != nil {
		return nil, err
	}
	outputCert.IssuerCertificate = string(issuerCertificate)

	return outputCert, nil
}
