# `gatewayexport`

```
gatewayexport <file>
```

## Synopsis

The `gatewayexport` tool exports content from a Gateway deployment.

See the [`gatewayimport`](gatewayimport.md) document for more information regarding [`gatewayimport`](gatewayimport.md), which provides the inverse “importing” capability.

Run `gatewayexport` from the system command line, not the Gateway shell.

## Arguments

## Options

### `--help`

Returns information on the options and use of **gatewayexport**.

### `--quiet`

Runs **gatewayexport** in a quiet mode that attempts to limit the amount of output.

### `--uri=<connectionString>`

Specify a resolvable URI connection string (enclose in quotes) to connect to the Gateway deployment.

```
--uri "http://host[:port]"
```

### `--format=<format>`

Format to use for the exported data (if can not be detected defaults to JSON-LD)

### `--out=<filename>`

Specifies the location and name of a file to export the data to. If you do not specify a file, **gatewayexport** writes data to the standard output (e.g. “stdout”).
