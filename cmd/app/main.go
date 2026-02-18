package main

import (
	"log"

	"github.com/kabirnayeem99/gohttplearn/internal/server"
)

func main() {

	srv := server.New(":9382")

	log.Println("Serving on", srv.Addr)

	log.Fatal(srv.Start())
}
