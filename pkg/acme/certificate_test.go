package acme

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCertificate(t *testing.T) {
	tests := []struct {
		name    string
		certPEM string
	}{
		{name: "can get certificate expire time", certPEM: "-----BEGIN CERTIFICATE-----\nMIIFFTCCA/2gAwIBAgITAPo8ZHN0QLKIpol7efRSWFaEJzANBgkqhkiG9w0BAQsF\nADBZMQswCQYDVQQGEwJVUzEgMB4GA1UEChMXKFNUQUdJTkcpIExldCdzIEVuY3J5\ncHQxKDAmBgNVBAMTHyhTVEFHSU5HKSBBcnRpZmljaWFsIEFwcmljb3QgUjMwHhcN\nMjMwOTAxMTEyNDAyWhcNMjMxMTMwMTEyNDAxWjAYMRYwFAYDVQQDEw1hbnNvbi5p\ndHN0LmNuMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAlJR27MQIKLGa\n0vNhLZLUCOJMBJnzkT4Ld19FC3Fb1HCJiyInay+EmF7OdHLUnPtQsDWCMnirs741\ntaiPKSy9S6qoHAcpJi6wLUeuO5w5Jc63x+URHVS30wX4rW+OxD4XhtYZ1C6xnnVE\n94Cy79rofu1o7w/7qBMuKj4BmM0Si2SJsgkV7Hw9cD/NRCnQ/FBclpG96l2tAod7\nBhMn9McL2aNc+ad0rPOjvjNkt9GZE3NP6q4sV5Y5G2lomYSHcHm3072YqY+ohlaG\nUhc9pcDV8HEcEyeTNTHD7/gdfH6P0uw363xClYgbAH09z8DTonR/NRga6cHx2U7U\n1HMOfJpH4QIDAQABo4ICFTCCAhEwDgYDVR0PAQH/BAQDAgWgMB0GA1UdJQQWMBQG\nCCsGAQUFBwMBBggrBgEFBQcDAjAMBgNVHRMBAf8EAjAAMB0GA1UdDgQWBBROOUOd\nMjyTDWOU7+rZMLyJgSharzAfBgNVHSMEGDAWgBTecnpI3zHDplDfn4Uj31c3S10u\nZTBdBggrBgEFBQcBAQRRME8wJQYIKwYBBQUHMAGGGWh0dHA6Ly9zdGctcjMuby5s\nZW5jci5vcmcwJgYIKwYBBQUHMAKGGmh0dHA6Ly9zdGctcjMuaS5sZW5jci5vcmcv\nMBgGA1UdEQQRMA+CDWFuc29uLml0c3QuY24wEwYDVR0gBAwwCjAIBgZngQwBAgEw\nggECBgorBgEEAdZ5AgQCBIHzBIHwAO4AdQDtq50d3YNzlZ/1Kojka7S8w8TMTXaK\nYMz/TjYtf7jWaAAAAYpQs06RAAAEAwBGMEQCIHFk0FDqDbRf92GYbhkyF8mduYp0\nOmxdkhQ11GOIq8AgAiB5OPoASB4qGBpMjgNFqkw/KCVLASVSw9YWzQCBSsKLJQB1\nALDMg+Wl+X1rr3wJzChJBIcqx+iLEyxjULfG/SbhbGx3AAABilCzUJYAAAQDAEYw\nRAIgekunlzNiw/312g0Xz5KCvPi9FUi9mdWJZHXuaBrjZHsCIE415dY3V44trnCv\nppZABN04lxwaaI43My2kzF9arcVsMA0GCSqGSIb3DQEBCwUAA4IBAQBmf4PgeHTK\nXbi9dCRi1tnrHyT/8Hh80vEmYcATYLiFU20VLp2ze1uQHhKhGsDvlJ58+MAoQaJg\nnK2mZRBjtyyJS0/1FEM1jR5O0yU4mlSZJ2jQ+nWtbHptP6F9YT8EAyVrQpjGJTvC\nwngUYGBzsAjdkk+MSrvOJyrxQMI0nhOq0T4gBja8vU9SYX/IsDBaDiaSE4ZXeC4M\nU5cUah8KCbDzaJtHSaWx0UPlC2SFhGhPyahcJr+d+PCViuckJX1hUgF6lymLcnJt\nPYST04onKZNrX855MB8kgfyAytYANdkfjv6W58d8Q02IvvwyceUm2dB6qdk3Xs9Y\nlKG52Ane+ar9\n-----END CERTIFICATE-----"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expire, err := GetCertificateExpireTime(tt.certPEM)

			if err != nil {
				t.Error(err)
			}

			now := time.Now()
			if now.Unix() > expire {
				assert.False(t, now.Unix() > expire)
			} else {
				assert.True(t, now.Unix() < expire)
			}
		})
	}
}
