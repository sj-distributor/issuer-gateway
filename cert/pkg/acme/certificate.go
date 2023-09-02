package acme

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
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
