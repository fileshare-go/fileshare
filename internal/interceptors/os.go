package interceptors

import (
	"context"

	"github.com/chanmaoganda/fileshare/internal/osinfo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func UnaryOSInfoInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {

	ctx = metadata.AppendToOutgoingContext(ctx, "os", osinfo.OsInfo())

	return invoker(ctx, method, req, reply, cc, opts...)
}

func StreamOSInfoInterceptor(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn,
	method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {

	ctx = metadata.AppendToOutgoingContext(ctx, "os", osinfo.OsInfo())

	return streamer(ctx, desc, cc, method, opts...)
}
