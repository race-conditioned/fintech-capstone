package grpc_transport

import (
	"context"
	"fintech-capstone/m/v2/internal/api_gateway/ports"

	"google.golang.org/grpc/metadata"
)

func metaFromGRPC(ctx context.Context, fullMethod string) ports.RequestMeta {
	first := func(vals []string) string {
		if len(vals) > 0 {
			return vals[0]
		}
		return ""
	}
	md, _ := metadata.FromIncomingContext(ctx)
	return ports.RequestMeta{
		ClientID:  first(md.Get("x-client-id")),
		RequestID: first(md.Get("x-request-id")),
		TraceID:   first(md.Get("x-trace-id")),
		Protocol:  "grpc",
		Target:    fullMethod,
	}
}
