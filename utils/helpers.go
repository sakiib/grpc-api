package utils

import (
	"github.com/sakiib/grpc-api/certs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func CertOption(tlsEnabled string) credentials.TransportCredentials {
	if tlsEnabled != "true" {
		return nil
	}
	crt := certs.GetServerCerts()
	if crt == nil {
		return nil
	}
	return credentials.NewServerTLSFromCert(crt)
}

func GetDialOption(tlsEnabled string) grpc.DialOption {
	if tlsEnabled != "true" {
		return grpc.WithInsecure()
	}
	return grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(certs.GetCertPool(), ""))
}
