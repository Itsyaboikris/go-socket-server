# Go TCP Socket Server

A high-performance TCP socket server implementation in Go with support for concurrent connections, graceful shutdown, and comprehensive testing.

## Features

- TCP socket server with concurrent connection handling
- Connection pooling and rate limiting
- Graceful shutdown support
- Configurable timeouts and connection limits
- Comprehensive test suite including load tests
- Client library with connection management
- Support for up to 1000 concurrent connections


## Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/go_socket_server.git

# Change into project directory
cd go_socket_server

# Install dependencies
go mod tidy

```

## Usage
### Starting the Server

``` bash
# Run the server
go run main.go
```

The server will start listening on port 9000 by default.

### Testing

``` bash
# Run all tests
go test -v ./test

# Run load tests
go test -v ./test -bench=BenchmarkConcurrentConnections -benchtime=10s

# Run specific test
go test -v ./test -run TestServerResponse
```
### Performance
The server is designed to handle up to 1000 concurrent connections. Load tests are included to verify performance under heavy load.

Example benchmark results:

```
BenchmarkConcurrentConnections-8    100    15234859 ns/op
```

### Contact
Your Name - @itsyaboikris
Project Link: https://github.com/itsyaboikris/go_socket_server


