# Migration

## From different backend

First you need to dump all the data from old backend \(`pq` extension is important\):

```bash
./gateway dump -d <backend> -a <address> -o ./data.pq.gz
```

or using config file:

```bash
./gateway dump -c <config> -o ./data.pq.gz
```

And load the data into a new backend and/or database:

```bash
./gateway load --init -d <new-backend> -a <new-address> -i ./data.pq.gz
```

or using config file:

```bash
./gateway load --init -c <new-config> -i ./data.pq.gz
```

### Dump via text format

An above guide uses Gateway-specific binary format to avoid encoding and parsing overhead and to compress output file better.

As an alternative, a standard nquads file format can be used to dump and load data \(note `nq` extension\):

```bash
./gateway dump -c <config> -o ./data.nq.gz
./gateway load --init -c <new-config> -i ./data.nq.gz
```

