package dt

import (
	"context"
	"net/http"
	"time"
)

// HTTPServer implements inbound.Server using net/http.
type HTTPServer struct {
	srv *http.Server
}

// NewHTTPServer creates a new HTTPServer with the given configuration.
func NewHTTPServer(
	addr string,
	handler http.Handler,
	readHeader,
	read,
	write,
	idle time.Duration,
) (*HTTPServer, error) {
	srv := &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: readHeader,
		ReadTimeout:       read,
		WriteTimeout:      write,
		IdleTimeout:       idle,
	}

	return &HTTPServer{srv: srv}, nil
}

// Start starts the HTTP server and listens for incoming requests.
func (s *HTTPServer) Start(ctx context.Context) error {
	errCh := make(chan error, 1)
	go func() {
		errCh <- s.srv.ListenAndServe()
	}()
	select {
	case <-ctx.Done():
		// caller will invoke Shutdown
		return ctx.Err()
	case err := <-errCh:
		if err == http.ErrServerClosed {
			return nil
		}
		return err
	}
}

// Shutdown gracefully shuts down the HTTP server.
func (s *HTTPServer) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
