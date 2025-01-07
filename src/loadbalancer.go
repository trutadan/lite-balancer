package src

import (
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

type LoadBalancer struct {
	servers 	 	[]string
	healthyServers 	map[string]bool
	currentIndex 	int
	mutex 		 	sync.Mutex
}

func NewLoadBalancer(servers []string) *LoadBalancer {
	healthyServers := make(map[string]bool)
	for _, server := range servers {
		healthyServers[server] = true
	}

	return &LoadBalancer{
		servers: servers,
		healthyServers: healthyServers,
		currentIndex: 0,
		mutex: sync.Mutex{},
	}
}

func (lb *LoadBalancer) NextServer() string {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()

	for i := 0; i < len(lb.servers); i++ {
		server := lb.servers[lb.currentIndex]
		lb.currentIndex = (lb.currentIndex + 1) % len(lb.servers)

		if lb.healthyServers[server] {
			return server
		}
	}

	return ""
}

func (lb *LoadBalancer) StartLoadBalancerServer(port string) {
	// create a listener
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer listener.Close()

	log.Printf("Load balancer started on port %s...", port)

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

			server := lb.NextServer()
			if server == "" {
				log.Printf("No healthy servers available...")
				conn.Close()
				return
			}

			handleConnection(conn, server)
		}(conn)
	}
}

func (lb *LoadBalancer) HealthCheck(interval time.Duration, checkURL string) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		for _, server := range lb.servers {
			go func(server string) {
				resp, err := http.Get(server + checkURL)
				lb.mutex.Lock()
				defer lb.mutex.Unlock()

				if err != nil || resp.StatusCode != http.StatusOK {
					log.Printf("Server %s is unhealthy", server)
					lb.healthyServers[server] = false
				} else {
					log.Printf("Server %s is healthy", server)
					lb.healthyServers[server] = true
				}

				if resp != nil {
					resp.Body.Close()
				}
			}(server)
		}
	}
}