package acme

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/pygzfei/issuer-gateway/utils"
	"time"
)

// GetCertificateExpireTime 获取证书过期时间
func GetCertificateExpireTime(certPEM string) (expire int64, err error) {

	// 解码 PEM 格式的证书
	block, _ := pem.Decode([]byte(certPEM))
	if block == nil {
		return time.Now().Add(-24 * time.Hour).Unix(), errors.New("error decoding certificate PEM")
	}

	// 解析证书
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return time.Now().Add(-24 * time.Hour).Unix(), errors.New(fmt.Sprintf("Error parsing certificate: %s", err))
	}

	// 获取证书的过期日期 (并且提前3天作为过期时间)
	return cert.NotAfter.Add(-72 * time.Hour).Unix(), nil
}

func EncryptCertificate(certInfo *certificate.Resource, secret string) (certificateEncrypt, privateKey, issuerCertificate string, expire int64, err error) {

	certificateEncrypt, err = utils.Encrypt(string(certInfo.Certificate), secret)

	if err != nil {
		return "", "", "", 0, err
	}

	privateKey, err = utils.Encrypt(string(certInfo.PrivateKey), secret)
	if err != nil {
		return "", "", "", 0, err
	}

	if string(certInfo.IssuerCertificate) != "" {
		issuerCertificate, err = utils.Encrypt(string(certInfo.IssuerCertificate), secret)
		if err != nil {
			return "", "", "", 0, err
		}
	}

	expire, err = GetCertificateExpireTime(string(certInfo.Certificate))
	if err != nil {
		return "", "", "", 0, err
	}

	return certificateEncrypt, privateKey, issuerCertificate, expire, nil
}

func DecryptCertificate(certificateEncrypt, privateKeyEncrypt, issuerCertificateEncrypt, secret string) (certificateDecrypt, privateKeyDecrypt, issuerCertificateDecrypt string, err error) {

	certificateDecrypt, err = utils.Decrypt(certificateEncrypt, secret)
	if err != nil {
		return "", "", "", err
	}

	privateKeyDecrypt, err = utils.Decrypt(privateKeyEncrypt, secret)
	if err != nil {
		return "", "", "", err
	}
	if issuerCertificateEncrypt != "" {
		issuerCertificateDecrypt, err = utils.Decrypt(issuerCertificateEncrypt, secret)
		if err != nil {
			return "", "", "", err
		}
	}

	return certificateDecrypt, privateKeyDecrypt, issuerCertificateDecrypt, nil
}
