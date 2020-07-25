# `gatewayimport`

```
gatewayimport <file>
```

## Synopsis

The `gatewayimport` tool imports content created by [`gatewayexport`](gatewayexport.md), or potentially, another third-party export tool.

See the [`gatewayexport`](gatewayexport.md) document for more information regarding [`gatewayexport`](gatewayexport.md), which provides the inverse “exporting” capability.

Run `gatewayimport` from the system command line, not the Gateway shell.

## Arguments

### `file`

Specifies the location and name of a file containing the data to import. If you do not specify a file, **gatewayimport** reads data from standard input (e.g. “stdin”).

## Options

### `--help`

Returns information on the options and use of **gatewayimport**.

### `--quiet`

Runs **gatewayimport** in a quiet mode that attempts to limit the amount of output.

### `--uri=<connectionString>`

Specify a resolvable URI connection string (enclose in quotes) to connect to the Gateway deployment.

```
--uri "http://host[:port]"
```

### `--format=<format>`

Format of the provided data (if can not be detected defaults to JSON-LD)
