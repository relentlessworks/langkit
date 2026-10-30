package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/relentlessworks/langkit/internal/api"
	"github.com/relentlessworks/langkit/internal/config"
)

func main() {
	cfg := config.Load()

	handler := api.NewHandler(cfg.NoAuth)
	mux := handler.Routes()

	server := &http.Server{
		Addr:         cfg.Addr,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	fmt.Printf("langkit — agentic-first language detection service\n")
	fmt.Printf("listening on %s\n", cfg.Addr)
	fmt.Printf("auth: %v\n", !cfg.NoAuth)
	fmt.Printf("GET /help for documentation\n\n")

	log.Fatal(server.ListenAndServe())
}
