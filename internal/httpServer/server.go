package httpServer

import (
	"context"
	"net"
	nethttp "net/http"
	"time"
)

type Server struct {
	srv *nethttp.Server
}

func NewServer(addr string, h nethttp.Handler) *Server {
	return &Server{
		srv: &nethttp.Server{
			Addr:              addr,
			Handler:           h,
			ReadTimeout:       15 * time.Second,
			ReadHeaderTimeout: 5 * time.Second,
			WriteTimeout:      15 * time.Second,
			IdleTimeout:       60 * time.Second,
			MaxHeaderBytes:    1 << 20,
		},
	}
}

func (s *Server) Listen(l net.Listener) error        { return s.srv.Serve(l) }
func (s *Server) Shutdown(ctx context.Context) error { return s.srv.Shutdown(ctx) }

func (s *Server) Addr() string { return s.srv.Addr }
