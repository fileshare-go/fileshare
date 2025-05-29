package interceptors

import (
	"context"
	"strings"

	"github.com/chanmaoganda/fileshare/internal/config"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

func BlackFilterInterceptor(settings *config.Settings) grpc.UnaryServerInterceptor {
	blockedIps := make(map[string]bool)
	for _, ip := range settings.BlockedIps {
		blockedIps[ip] = true
	}

	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		peer, ok := peer.FromContext(ctx)
		if !ok {
			return nil, status.Errorf(status.Code(err), "Failed to get client addr: %v", err)
		}
		parts := strings.Split(peer.Addr.String(), ":")
		addr := parts[0]
		if blockedIps[addr] {
			logrus.Debugf("Ip from %s is blocked", addr)
			return nil, status.Errorf(codes.Aborted, "Your ip addr is blocked")
		}

		return handler(ctx, req)
	}
}
