package interceptors

import (
	"context"

	"google.golang.org/grpc"
)

func UnaryRecorderInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {

	// record := model.Record{

	// }

	// switch req.(type) {

	// }

	// service.Orm().Create()

	return invoker(ctx, method, req, reply, cc, opts...)
}
