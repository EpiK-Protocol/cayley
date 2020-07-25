# Install Gateway

## Install Gateway on Ubuntu

```text
snap install --edge --devmode gateway
```

## Install Gateway on macOS

### Install Homebrew

macOS does not include the Homebrew brew package by default. Install brew using the [official instructions](https://brew.sh/#install)

### Install Gateway

```bash
brew install gateway
```

## Install Gateway with Docker

```bash
docker run -p 64210:64210 epik-protocol/gateway
```

For more information see [Container Documentation](deployment/container.md)

## Build from Source

See instructions in [Contributing](getting-involved/contributing.md)

