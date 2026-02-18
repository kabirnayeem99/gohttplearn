package server

func (s *Server) routes() {
	s.mux.HandleFunc("/", handleHello)
	s.mux.HandleFunc("/goodbye", handleGoodBye)
}
