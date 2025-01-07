package main

import (
	"lite-balancer/src"
	"time"
)

func main() {
	loadBalancerPort := "8080"

	lb := src.NewLoadBalancer([]string{"http://localhost:8081", "http://localhost:8082", "http://localhost:8083"})

	go lb.HealthCheck(10*time.Second, "/health")

	lb.StartLoadBalancerServer(loadBalancerPort)
}
