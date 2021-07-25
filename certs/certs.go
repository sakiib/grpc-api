package certs

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
)

var (
	// Cert is the self-signed certificate
	Cert tls.Certificate
	// CertPool contains the self-signed certificate
	CertPool *x509.CertPool
)

func PrepareCerts() {
	crt, err := ioutil.ReadFile("certs/server-cert.pem")
	if err != nil {
		log.Fatalf("failed to read certs/server-cert.pem, %s", err.Error())
	}

	key, err := ioutil.ReadFile("certs/server-key.pem")
	if err != nil {
		log.Fatalf("failed to read certs/server-key.pem, %s", err.Error())
	}

	Cert, err = tls.X509KeyPair(crt, key)
	if err != nil {
		log.Fatalf("failed to create certs key pair, %s", err.Error())

	}

	Cert.Leaf, err = x509.ParseCertificate(Cert.Certificate[0])
	if err != nil {
		log.Fatalf("failed to parse certificate, %s", err.Error())
	}

	CertPool = x509.NewCertPool()
	CertPool.AddCert(Cert.Leaf)
}
