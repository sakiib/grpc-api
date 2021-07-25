package utils

import (
	"crypto/tls"
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
