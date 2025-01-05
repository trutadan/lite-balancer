package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
)

func main() {
	loadBalancerPort := "8080"
	backendServerIP := "http://localhost:8081"

	// create a listener
	listener, err := net.Listen("tcp", ":"+loadBalancerPort)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer listener.Close()

	log.Printf("Load balancer started on port %s", loadBalancerPort)

	// wait group to wait for all goroutines to finish
	var wg sync.WaitGroup

	for {
		// accept a connection
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}

		// increment the wait group counter
		wg.Add(1)

		// handle the connection
		go func(conn net.Conn) {
			defer wg.Done()
			handleConnection(conn, backendServerIP)
		}(conn)
	}
}

func handleConnection(clientConnection net.Conn, backendServerIP string) {
	defer clientConnection.Close()

	// parse the request from the client
	req, err := http.ReadRequest(bufio.NewReader(clientConnection))
	if err != nil {
		log.Printf("Failed to read request: %v", err)
		return
	}

	// log the incoming request
	log.Printf("Received request: %s %s", req.Method, req.URL)

	// forward the request to the backend server
	resp, err := forwardRequest(req, backendServerIP)
	if err != nil {
		log.Printf("Failed to forward request: %v", err)
		return
	}

	// write the response back to the client
	err = resp.Write(clientConnection)
	if err != nil {
		log.Printf("Failed to write response to client: %v", err)
	}
}

func forwardRequest(req *http.Request, backendServerIP string) (*http.Response, error) {
	// create a new request for the backend server
	proxyReq, err := http.NewRequest(req.Method, backendServerIP+req.URL.Path, req.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for backend: %w", err)
	}

	// copy headers from the original request
	proxyReq.Header = req.Header

	// send the request to the backend server
	client := &http.Client{}
	resp, err := client.Do(proxyReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request to backend: %w", err)
	}

	return resp, nil
}