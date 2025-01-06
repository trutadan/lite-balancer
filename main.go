package main

import (
	"lite-balancer/src"
)

func main() {
	loadBalancerPort := "8080"

	lb := src.NewLoadBalancer([]string{"http://localhost:8081", "http://localhost:8082"})

	lb.StartLoadBalancerServer(loadBalancerPort)
}
