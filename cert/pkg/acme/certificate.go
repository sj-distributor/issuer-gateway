package acme

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"time"
)

func Certificate() {
	// 读取证书文件
	certPEM, err := ioutil.ReadFile("your_certificate.pem")
	if err != nil {
		fmt.Println("Error reading certificate file:", err)
		return
	}

	// 解码 PEM 格式的证书
	block, _ := pem.Decode(certPEM)
	if block == nil {
		fmt.Println("Error decoding certificate PEM")
		return
	}

	// 解析证书
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		fmt.Println("Error parsing certificate:", err)
		return
	}

	// 获取证书的过期日期
	expiration := cert.NotAfter
	fmt.Println("Certificate Expires On:", expiration.Format(time.RFC3339))
}
