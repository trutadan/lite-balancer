package src

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/http"
)

func handleConnection(clientConnection net.Conn, backendServerIP string) {
	defer clientConnection.Close()

	// parse the request from the client
	req, err := http.ReadRequest(bufio.NewReader(clientConnection))
	if err != nil {
		log.Printf("Failed to read request: %v", err)
		return
	}

	// log the incoming request
	log.Printf("Received request of type %s to %s, redirected it to %s", req.Method, req.URL, backendServerIP)

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