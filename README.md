# gRPC Blog Service

A high-performance Blog Management API built with Go and gRPC. This service demonstrates enterprise patterns: interface-based design, thread-safe in-memory persistence, and automated mock-driven testing.

### Key Features

- **gRPC API** - Full CRUD implementation using Protobuf and modern gRPC-Go
- **Concurrency** - Thread-safe store utilizing sync.RWMutex (optimized for multiple readers).
- **Configuration** - Centralized management via Viper (YAML/Environment variable support).
- **Quality Assurance** - Unit tests with Uber GoMock and filtered coverage reporting.
- **Automation** - Integrated with Buf and a robust Makefile for consistent developer workflows.

## Technology Stack

| Category | Tool | Purpose |
|-----------|---------|---------|
| Language | Go 1.24.5 | Primary runtime |
| RPC | gRPC / Protobuf | API definition and serialization |
| Config | Viper | Configuration management |
| Mocking | Uber GoMock | Interface-based unit testing |
| Tooling | Buf / Make | Code generation and automation |
 

## Project Structure

```
grpc-blog-service/
├── proto/                    # Protocol buffer definitions
│   ├── blog.proto           # Service and message definitions
│   ├── blog.pb.go           # Generated Go code
│   └── blog_grpc.pb.go      # Generated gRPC code
├── server/                  # gRPC server implementation
│   ├── main.go             # Server entry point
│   └── server.go           # Service handler implementations
├── client/                  # gRPC client for testing
│   ├── main.go             # Client entry point
│   └── handlers.go         # Test request handlers
├── internal/               # Internal utilities and packages
│   ├── store/             # Data store implementation
│   │   ├── store.go       # Store interface and implementation
│   │   └── store_test.go  # Unit tests for store
│   └── mocks/             # Mock implementations
│       └── mock_store.go  # Generated mock store
├── config/                # Configuration
│   ├── config.go          # Configuration loader
│   └── config.yaml        # Configuration file
├── buf.gen.yaml           # Buf code generation config
├── go.mod                 # Go module definition
├── Makefile              # Build automation
└── README.md             # This file
```



## Installation & Setup

### Prerequisites

- Go 1.24.5 or higher
- Make (optional, for running make commands)

### Clone the Repository

```bash
git clone https://github.com/Chetas1/grpc-blog-service.git
cd grpc-blog-service
```

### Install Dependencies

```bash
go mod download
go mod tidy
```



## Running the Service

### Start the Server

```bash
make server
# or
go run ./server
```

The server will start on `localhost:9090` by default.

### Run the Client (Test)

In a separate terminal:

```bash
make client
# or
go run ./client
```

### Configuration

The service uses YAML configuration. Edit `config/config.yaml` to customize:

```yaml
GrpcServer:
  Host: localhost
  Port: 9090
  Protocol: tcp

GrpcClient: 
  ServerAddress: "localhost:9090"
```

## Testing

### Run Unit Tests

```bash
make test
# or
go test ./... -v
```

### Run Tests with Coverage

```bash
make test
```

This generates a `coverage.out` file and displays coverage statistics.