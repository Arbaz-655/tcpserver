package main

import (
	"log"
	"os"

	"github.com/Arbaz/tcp/config"
	"github.com/Arbaz/tcp/server"
)

func main() {
	// Create or open the log file
	logFile, err := os.OpenFile("server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()

	// Set log output to the log file
	log.SetOutput(logFile)

	// Load configuration
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Define TCP server with provided Config
	server := server.NewServer(cfg.ServerConfig)

	//Use Go routine to start Server
	go func() {
		if err := server.Start(); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
		log.Println("Server Start Called")
	}()

	//keep the main function running
	log.Println("Main: Waiting for server to run...")
	select {}
}
