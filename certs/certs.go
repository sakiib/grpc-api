package certs

import (
	"github.com/sakiib/grpc-api/utils"
	"google.golang.org/grpc/credentials"
)

func CertOption(tlsEnabled string) credentials.TransportCredentials {
	if tlsEnabled != "true" {
		return nil
	}
	crt := utils.GetServerCerts()
	if crt == nil {
		return nil
	}
	return credentials.NewServerTLSFromCert(crt)
}
