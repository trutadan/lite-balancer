# lite-balancer
This project implements a lightweight, application-layer HTTP Load Balancer in Go.<br />
The load balancer distributes client requests to backend servers using a round-robin scheduling algorithm. It handles multiple concurrent connections efficiently with goroutines.<br />
It also periodically performs HTTP GET requests to backend servers' health check route. Marks servers as unhealthy if they fail to respond. Automatically reintroduces servers when they recover.

### Running the code for the Load Balancer:
1. #### Clone this repository
```
git clone https://github.com/trutadan/lite-balancer.git
cd lite-balancer
```

2. #### Start the backend servers
Use the provided server.py script.<br />
Start two or more instances on different ports:
```
python server.py 8081
python server.py 8082
```

3. #### Start the load balancer
```
go run main.go
```

4. #### Test the load balancer
Send requests to the load balancer:
```
curl http://localhost:8080/
```

The responses should alternate between the backend servers.

5. #### Simulate server failure
Stop one of the backend servers (e.g., 8081) and observe that requests are routed only to healthy servers.<br />
Restart the server to verify that it re-enters the pool.