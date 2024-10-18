package zrpc

import (
	"context"
	"github.com/apache/skywalking-go/plugins/core/operator"
	"github.com/apache/skywalking-go/plugins/core/tracing"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type ServerMiddlewareInterceptor struct {
}

// BeforeInvoke intercepts the HTTP request before invoking the handler.
func (h *ServerMiddlewareInterceptor) BeforeInvoke(invocation operator.Invocation) error {
	server := invocation.CallerInstance().(*zrpc.RpcServer)
	server.AddUnaryInterceptors(RpcServeInterceptor(invocation))
	return nil
}

// AfterInvoke processes after the HTTP request has been handled.
func (h *ServerMiddlewareInterceptor) AfterInvoke(invocation operator.Invocation, result ...interface{}) error {
	return nil
}

// RpcServeInterceptor is a grpc server interceptor that creates a new span for each incoming request.
var RpcServeInterceptor = func(invocation operator.Invocation) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		span, err := tracing.CreateEntrySpan(info.FullMethod, func(headerKey string) (string, error) {
			md, ok := metadata.FromIncomingContext(ctx)
			if !ok {
				return "", nil
			}
			values := md.Get(headerKey)
			if len(values) == 0 || len(values[0]) == 0 {
				return "", nil
			}
			return values[0], nil
		}, tracing.WithComponent(5023),
			tracing.WithLayer(tracing.SpanLayerRPCFramework),
			tracing.WithTag("transport", "gRPC"))
		if err != nil {
			return handler(ctx, req)
		}
		defer span.End()

		reply, err := handler(ctx, req)
		if err != nil {
			span.Error(err.Error())
		}
		return reply, err
	}
}
