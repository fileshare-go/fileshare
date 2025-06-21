package fileshare

import (
	"github.com/chanmaoganda/fileshare/internal/config"
	"github.com/chanmaoganda/fileshare/internal/pkg/certs"
	"github.com/chanmaoganda/fileshare/internal/pkg/interceptors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewCredentialClientConn(settings *config.Config) (*grpc.ClientConn, error) {
	credentialLoader, err := certs.NewCredentialLoader(settings.CertsPath)
	if err != nil {
		return nil, err
	}

	return grpc.NewClient(
		settings.GrpcAddress,
		grpc.WithTransportCredentials(credentialLoader.ServerCredentials),
	)
}

func NewCredentialServerConn(settings *config.Config) (*grpc.Server, error) {
	credentialLoader, err := certs.NewCredentialLoader(settings.CertsPath)
	if err != nil {
		return nil, err
	}
	grpcServer := grpc.NewServer(grpc.Creds(credentialLoader.ServerCredentials))
	return grpcServer, nil
}

func NewClientConn(settings *config.Config) (*grpc.ClientConn, error) {
	return grpc.NewClient(settings.GrpcAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(interceptors.UnaryOSInfoInterceptor),
		grpc.WithStreamInterceptor(interceptors.StreamOSInfoInterceptor),
	)
}

func NewServerConn(settings *config.Config) (*grpc.Server, error) {
	return grpc.NewServer(
		grpc.UnaryInterceptor(interceptors.BlackFilterInterceptor(settings)),
	), nil
}
