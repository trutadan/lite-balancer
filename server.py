from http.server import SimpleHTTPRequestHandler, HTTPServer
import sys

class HealthCheckHandler(SimpleHTTPRequestHandler):
    def do_GET(self):
        if self.path == "/health":
            self.send_response(200)
            self.end_headers()
            self.wfile.write(b"OK")
        else:
            super().do_GET()

if __name__ == "__main__":
    # default port
    port = 8081
    
    # check if a port is provided as a command-line argument
    if len(sys.argv) > 1:
        try:
            port = int(sys.argv[1])
        except ValueError:
            print("Invalid port number. Using default port 8081.")
    
    httpd = HTTPServer(("localhost", port), HealthCheckHandler)
    print(f"Serving on port {port}...")
    httpd.serve_forever()
