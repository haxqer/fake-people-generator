package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"fake-people-generator/internal/api"
	"fake-people-generator/internal/generator"
)

func main() {
	addr := strings.TrimSpace(os.Getenv("ADDR"))
	if addr == "" {
		port := strings.TrimSpace(os.Getenv("PORT"))
		if port != "" {
			if strings.HasPrefix(port, ":") {
				addr = port
			} else {
				addr = ":" + port
			}
		}
	}
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
