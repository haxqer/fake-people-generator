package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"fake-people-generator/internal/api"
	"fake-people-generator/internal/generator"
)

func main() {
	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":8080"
	}

	srv := &http.Server{
		Addr:              addr,
		Handler:           api.NewServer(generator.New(nil)),
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("fake people generator listening on %s", addr)

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
