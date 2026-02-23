package server

func (s *Server) routes() {
	s.mux.HandleFunc("/{$}", handleHello)
	s.mux.HandleFunc("/goodbye", handleGoodBye)
	s.mux.HandleFunc("/hello", handleHelloParameterized)
	s.mux.HandleFunc("/greetings/hello", handleGreetingsHello)
	s.mux.HandleFunc("/greetings/{user}/hello", handlerGreetingsUserHello)
}
