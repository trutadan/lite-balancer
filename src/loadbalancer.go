package src

import (
	"log"
	"net"
	"sync"
)

type LoadBalancer struct {
	servers 	 []string
	currentIndex int
	mutex 		 sync.Mutex
}

func NewLoadBalancer(servers []string) *LoadBalancer {
	return &LoadBalancer{
		servers: servers,
		currentIndex: 0,
		mutex: sync.Mutex{},
	}
}

func (lb *LoadBalancer) NextServer() string {
	lb.mutex.Lock()

	server := lb.servers[lb.currentIndex]
	lb.currentIndex = (lb.currentIndex + 1) % len(lb.servers)

	lb.mutex.Unlock()

	return server
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
			handleConnection(conn, lb.NextServer())
		}(conn)
	}
}