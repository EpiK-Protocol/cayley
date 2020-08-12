# Contributing

## Community Involvement

Join our community on [community.epik-protocol.io](https://) or other [Locations](locations.md).

## Simply building Gateway

If your version of Go &lt; 1.13, you need to run:

```text
export GO111MODULE=on
```

Follow the instructions for running Gateway locally:

```text
# clone project
git clone https://github.com/epik-protocol/epik-gateway-backend
cd gateway

# Download dependencies
go mod download

# Download web files (optional)

go run cmd/download_ui/download_ui.go

# Install packr 2

go get -u github.com/gobuffalo/packr/v2/packr2
```

# Generate static files go modules

packr2

# build the binary

go build ./cmd/gateway

# try the generated binary

```bash
./gateway help
```

Give it a quick test with:

```text
./gateway repl -i data/testdata.nq
```

To run the web frontend, replace the "repl" command with "http"

```text
./gateway http -i data/testdata.nq
```

You can now open the WebUI in your browser: [http://127.0.0.1:64210](http://127.0.0.1:64210)

## Hacking on Gateway

First, you'll need Go [\(version 1.11.x or greater\)](https://golang.org/doc/install) and a Go workspace. This is outlined by the Go team at [http://golang.org/doc/code.html](http://golang.org/doc/code.html) and is sort of the official way of going about it.

If your version of Go &lt; 1.13, you need to run:

```text
export GO111MODULE=on
```

If you just want to build Gateway and check out the source, or use it as a library, a simple `go get github.com/epik-protocol/epik-gateway-backend` will work!

But suppose you want to contribute back on your own fork \(and pull requests are welcome!\). A good way to do this is to set up your \$GOPATH and then...

```text
mkdir -p $GOPATH/src/github.com/epik-protocol
cd $GOPATH/src/github.com/epik-protocol
git clone https://github.com/$GITHUBUSERNAME/gateway
```

...where \$GITHUBUSERNAME is, well, your GitHub username :\) You'll probably want to add

```text
cd gateway
git remote add upstream http://github.com/epik-protocol/epik-gateway-backend
```

So that you can keep up with the latest changes by periodically running

```text
git pull --rebase upstream
```

With that in place, that folder will reflect your local fork, be able to take changes from the official fork, and build in the Go style.

For iterating, it can be helpful to, from the directory, run

```text
packr2 && go build ./cmd/gateway && ./gateway <subcommand> <your options>
```

Which will also resolve the relevant static content paths for serving HTTP.

**Reminder:** add yourself to CONTRIBUTORS and AUTHORS.

## Running Unit Tests

If your version of Go &lt; 1.13, you need to run:

```text
export GO111MODULE=on
```

First, `cd` into the `gateway` project folder and run:

```text
packr && go test ./...
```

If you have a Docker installed, you can also run tests for remote backend implementations:

```text
go test -tags docker ./...
```

If you have a Docker installed, you only want to run tests for a specific backend implementations eg. mongodb

```text
go test -tags docker ./graph/nosql/mongo
```

Integration tests can be enabled with environment variable:

```text
RUN_INTEGRATION=true go test ./...
```
