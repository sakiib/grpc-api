package certs

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
)

func GetServerCerts() *tls.Certificate {
	crt, err := ioutil.ReadFile("certs/server-cert.pem")
	if err != nil {
		return nil
	}
	key, err := ioutil.ReadFile("certs/server-key.pem")
	if err != nil {
		return nil
	}
	cert, err := tls.X509KeyPair(crt, key)
	if err != nil {
		return nil
	}
	return &cert
}

func GetCertPool() *x509.CertPool {
	Cert := GetServerCerts()
	Cert.Leaf, _ = x509.ParseCertificate(Cert.Certificate[0])

	CertPool := x509.NewCertPool()
	CertPool.AddCert(Cert.Leaf)

	return CertPool
}
