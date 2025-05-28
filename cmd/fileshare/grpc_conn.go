package fileshare

import (
	"github.com/chanmaoganda/fileshare/internal/config"
	"github.com/chanmaoganda/fileshare/internal/pkg/certs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewCredentialClientConn(settings *config.Settings) (*grpc.ClientConn, error) {
	credentialLoader, err := certs.NewCredentialLoader(settings.CertsPath)
	if err != nil {
		return nil, err
	}

	return grpc.NewClient(settings.GrpcAddress, grpc.WithTransportCredentials(credentialLoader.ServerCredentials))
}

func NewCredentialServerConn(settings *config.Settings) (*grpc.Server, error) {
	credentialLoader, err := certs.NewCredentialLoader(settings.CertsPath)
	if err != nil {
		return nil, err
	}
	grpcServer := grpc.NewServer(grpc.Creds(credentialLoader.ServerCredentials))
	return grpcServer, nil
}

func NewClientConn(settings *config.Settings) (*grpc.ClientConn, error) {
	return grpc.NewClient(settings.GrpcAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
}

func NewServerConn(settings *config.Settings) (*grpc.Server, error) {
	return grpc.NewServer(), nil
}
