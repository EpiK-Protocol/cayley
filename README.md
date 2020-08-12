# Contributing
If your version of Go < 1.13, you need to run:
```bash
export GO111MODULE=on
```

```bash
# clone project
git clone https://github.com/EpiK-Protocol/epik-gateway-backend.git
cd gateway

# Download dependencies
go mod download

# Install packr 2
go get -u github.com/gobuffalo/packr/v2/packr2
```

Run Web Demo
```bash
# Download web files
go run cmd/download_ui_demo/download_ui_demo.go

# Generate static files go modules
packr2

# Build the binary
go build ./cmd/gateway

# Try the binary
./gateway help

# Quick test
./gateway repl -i data/testdata.nq

# Run the web frontend
./gateway http -i data/demo/marvel_demo.nq
```