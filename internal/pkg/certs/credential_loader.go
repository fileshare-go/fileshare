package certs

import (
	"strings"

	"google.golang.org/grpc/credentials"
)

type CredentialLoader struct {
	ServerCredentials credentials.TransportCredentials
	ClientCredentials credentials.TransportCredentials
}

func NewCredentialLoader(certs_path string) (*CredentialLoader, error) {
	serverCert := strings.Join([]string{certs_path, "server.crt"}, "/")
	serverKey := strings.Join([]string{certs_path, "server.key"}, "/")

	server, err := credentials.NewServerTLSFromFile(serverCert, serverKey)
	if err != nil {
		return nil, err
	}

	clientCert := strings.Join([]string{certs_path, "ca.crt"}, "/")
	clientKey := strings.Join([]string{certs_path, "ca.key"}, "/")

	client, err := credentials.NewServerTLSFromFile(clientCert, clientKey)
	if err != nil {
		return nil, err
	}

	loader := &CredentialLoader{
		ServerCredentials: server,
		ClientCredentials: client,
	}

	return loader, nil
}
