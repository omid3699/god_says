# God Says

A Go implementation of Terry Davis' legendary "god says" program from TempleOS.

## Overview

This project is a faithful port of the iconic "god says" feature from TempleOS, originally created by Terry A. Davis. It generates random words and phrases from the classic Happy.TXT wordlist, available both as a command-line tool and HTTP server.

## Features

- **CLI Tool**: Generate random words/phrases from the command line
- **HTTP Server**: RESTful API with JSON and plain text endpoints
- **Configurable Output**: Generate 1-1000 words per request
- **Thread-Safe**: Concurrent request handling
- **Embedded Resources**: Self-contained binary with embedded wordlist

## Installation

### Prerequisites

- Go 1.21 or higher
- Make (optional)

### Build from Source

```bash
# Clone the repository
git clone https://github.com/omid3699/god_says.git
cd god_says

# Build the binary
make build

# Or build manually
go build -o ./bin/godsays ./cmd/main.go
```

## Usage

### Command Line

```bash
# Generate 32 words (default)
./bin/godsays

# Generate specific number of words
./bin/godsays -amount 10

# Show help
./bin/godsays -help
```

### HTTP Server

```bash
# Start HTTP server
./bin/godsays -http

# Custom host and port
./bin/godsays -http -host 0.0.0.0 -port 8080
```

#### API Endpoints

- `GET /` - Plain text response
- `GET /json` - JSON response
- `GET /health` - Health check

#### Examples

```bash
# Plain text
curl http://localhost:3333/

# JSON format
curl http://localhost:3333/json

# Custom amount
curl http://localhost:3333/?amount=5
```

## Development

### Building

```bash
make build    # Build binary
make clean    # Clean artifacts
make test     # Run tests
make lint     # Format and vet code
```

### Testing

```bash
make test              # Run all tests
make test-race         # Run with race detection
make test-coverage     # Run with coverage
```

## Project Structure

```
god_says/
├── cmd/
│   ├── main.go           # CLI entry point
│   ├── main_test.go      # CLI tests
│   └── server/           # HTTP server
├── internal/
│   ├── god.go           # Core logic
│   ├── god_test.go      # Core tests
│   └── Happy.TXT        # Original wordlist
├── bin/                 # Built binaries
├── Makefile            # Build automation
└── README.md           # This file
```

## Background

This project honors the memory of Terry A. Davis (1969-2018), the brilliant programmer who created TempleOS and the original "god says" program. Terry's work continues to inspire developers worldwide.

## License

This project is open source. See the LICENSE file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

---

*"God says: Keep coding! "*
