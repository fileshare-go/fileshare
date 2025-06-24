package util

import (
	"context"
	"encoding/base64"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

func PeerOs(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "unknown"
	}

	if osInfo, ok := md["os"]; ok && len(osInfo) != 0 {
		data, _ := base64.StdEncoding.DecodeString(osInfo[0])
		return string(data)
	}
	return "unknown"
}

func PeerAddress(ctx context.Context) string {
	peer, ok := peer.FromContext(ctx)
	if ok {
		return peer.Addr.String()
	}
	return "unknown"
}
