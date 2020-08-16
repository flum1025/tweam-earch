package main

import (
	"flag"
	"log"

	"github.com/flum1025/tweam-earch/internal/config"
	"github.com/flum1025/tweam-earch/internal/server"
)

func main() {
	configPath := flag.String("config", "", "config path")
	port := flag.Int("port", 3000, "port")
	flag.Parse()

	config, err := config.NewConfig(*configPath)
	if err != nil {
		log.Fatalf("failed to get config: %v", err)
	}

	srv, err := server.NewServer(config)
	if err != nil {
		log.Fatalf("failed to initialize server: %v", err)
	}

	if err := srv.Run(*port); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
