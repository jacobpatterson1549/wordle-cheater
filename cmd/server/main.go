package main

import (
	"log"
	"net"
	"net/http"
	"os"

	"github.com/jacobpatterson1549/wordle-cheater/cmd/server/config"
	"github.com/jacobpatterson1549/wordle-cheater/internal/server"
)

func main() {
	cfg, err := config.New(os.Stdout, os.LookupEnv, os.Args...)
	if err != nil {
		log.Fatalf("parsing configuration: %v", err)
	}

	addr := net.JoinHostPort(cfg.Host, cfg.Port)
	log.Println("Serving resume site at http://127.0.0.1" + addr)
	log.Println("Press Ctrl-C to stop")
	h := http.HandlerFunc(server.ServeHTTP)
	http.ListenAndServe(addr, h)
}
