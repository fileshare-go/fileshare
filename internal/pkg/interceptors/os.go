package interceptors

import (
	"context"

	"github.com/chanmaoganda/fileshare/internal/pkg/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func UnaryOSInfoInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		ctx = metadata.AppendToOutgoingContext(ctx, "os", util.OsInfo())
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func StreamOSInfoInterceptor() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		ctx = metadata.AppendToOutgoingContext(ctx, "os", util.OsInfo())
		return streamer(ctx, desc, cc, method, opts...)
	}
}
