package main

import (
	"log"
	"os"

	"poke/backend/internal/router"
)

func main() {
	r := router.New()

	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	if err := r.Run(addr); err != nil {
		log.Fatalf("server start failed: %v", err)
	}
}
