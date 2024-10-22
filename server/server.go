package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"
	"sync/atomic"

	"github.com/Arbaz/tcp/config"
)

/*
Create a TCP server that listens on 10 separate ports, processes the data received,
and responds back to each request with a unique identifier
*/

type Server struct {
	cfg       config.ServerConfig
	idCounter uint64
}

type Request struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Response struct {
	ID    uint64 `json:"id"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

/*
The NewServer function initializes a Server instance with the provided configuration.
*/
func NewServer(cfg config.ServerConfig) *Server {
	return &Server{
		cfg: cfg,
	}
}

/*
The Start method creates a WaitGroup to manage the goroutines for each port.
It iterates over the configured ports, starts a new goroutine for each port,
and waits for all goroutines to finish before returning.
*/
func (s *Server) Start() error {
	var wg sync.WaitGroup

	//Start server for each Port in different go routine
	for _, port := range s.cfg.Ports {
		wg.Add(1)

		go func(p int) {
			defer wg.Done()
			if err := s.listenOnPort(p); err != nil {
				log.Printf("Error listening on port %d: %v", p, err)
			}
		}(port)
	}

	wg.Wait()
	return nil
}

/*
The listenOnPort method sets up a TCP listener on a given port.
It listens for incoming connections and handles them in a separate goroutine.
*/
func (s *Server) listenOnPort(port int) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("failed to listen on port %d: %v", port, err)
	}
	defer listener.Close()

	log.Printf("Listening on port %d", port)

	//Infinite Loop for acceptin connecting till server is UP
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection on port %d: %v", port, err)
			continue
		}

		//Handling in another go routine so it doesn't block the accepting go routine
		go s.handleConnection(conn)
	}
}

/*
The handleConnection method processes incoming data from a connection.
It reads the data, logs it, and sends a response back to the client.
*/
func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		var req Request
		if err := json.Unmarshal(scanner.Bytes(), &req); err != nil {
			log.Printf("Error unmarshaling request: %v", err)
			continue
		}

		response := s.processRequest(req, conn)

		responseJSON, err := json.Marshal(response)
		if err != nil {
			log.Printf("Error marshaling response: %v", err)
			continue
		}

		if _, err := conn.Write(append(responseJSON, '\n')); err != nil {
			log.Printf("Error writing response: %v", err)
			return
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error reading from connection: %v", err)
	}
}

// We are returning uniqye response id
func (s *Server) processRequest(req Request, conn net.Conn) Response {

	// Generate a unique identifier
	id := atomic.AddUint64(&s.idCounter, 1)

	log.Printf("Request Received %v on connection %v", req.Key, conn.LocalAddr())

	// For now, we're just echoing back the received data with an ID
	return Response{
		ID:    id,
		Key:   req.Key,
		Value: req.Value,
	}
}
