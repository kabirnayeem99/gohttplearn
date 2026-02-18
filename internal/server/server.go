package server

import "net/http"

type Server struct {
	Addr string
	mux  *http.ServeMux
}

func New(addr string) *Server {
	mux := http.NewServeMux()
	s := &Server{
		Addr: addr,
		mux:  mux,
	}

	s.routes()

	return s
}

func (s *Server) Start() error {
	return http.ListenAndServe(s.Addr, s.mux)
}
