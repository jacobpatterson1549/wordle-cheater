package main

import (
	"log"
	"net"
	"net/http"
	"os"

	words "github.com/jacobpatterson1549/wordle-cheater"
	"github.com/jacobpatterson1549/wordle-cheater/cmd/server/config"
	"github.com/jacobpatterson1549/wordle-cheater/internal/server"
)

func main() {
	cfg, err := config.New(os.Stdout, os.LookupEnv, os.Args...)
	if err != nil {
		log.Fatalf("parsing configuration: %v", err)
	}

	h := server.NewHandler(words.WordsTextFile)
	addr := net.JoinHostPort(cfg.Host, cfg.Port)
	log.Println("Serving resume site at http://127.0.0.1" + addr)
	log.Println("Press Ctrl-C to stop")
	http.ListenAndServe(addr, h)
}
