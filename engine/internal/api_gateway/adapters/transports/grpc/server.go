package grpc_transport

import (
	"context"
	"net"

	"google.golang.org/grpc"
)

type GRPCServer struct {
	srv *grpc.Server
	ln  net.Listener
}

func NewGRPCServer(addr string, register func(*grpc.Server)) (*GRPCServer, error) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	gs := grpc.NewServer() // add interceptors later
	register(gs)           // register generated services here
	return &GRPCServer{srv: gs, ln: ln}, nil
}

func (s *GRPCServer) Start(ctx context.Context) error {
	errCh := make(chan error, 1)
	go func() { errCh <- s.srv.Serve(s.ln) }()
	select {
	case <-ctx.Done():
		s.srv.GracefulStop()
		return ctx.Err()
	case err := <-errCh:
		return err
	}
}

func (s *GRPCServer) Shutdown(ctx context.Context) error {
	done := make(chan struct{})

	go func() {
		s.srv.GracefulStop()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		s.srv.Stop()
		<-done // ensure cleanup before returning
		return ctx.Err()
	}
}
