# Container

## Running in Kubernetes

To run Gateway in K8S check [this docs section](k8s/k8s.md).

## Running in a container

A container exposing the HTTP API of Gateway is available.

### Running with default configuration

Container is configured to use BoltDB as a backend by default.

```text
docker run -p 64210:64210 -d epik-protocol/gateway:v0.1.0
```

New database will be available at [http://localhost:64210](http://localhost:64210).

### Custom configuration

To run the container one must first setup a data directory that contains the configuration file and optionally contains persistent files \(i.e. a boltdb database file\).

```text
mkdir data
cp gateway_example.yml data/gateway.yml
cp data/testdata.nq data/my_data.nq
# initialize and serve database
docker run -v $PWD/data:/data -p 64210:64210 -d epik-protocol/gateway:v0.1.0 -c /data/gateway.yml --init -i /data/my_data.nq
# serve existing database
docker run -v $PWD/data:/data -p 64210:64210 -d epik-protocol/gateway:v0.1.0 -c /data/gateway.yml
```

### Other commands

Container runs `gateway http` command by default. To run any other Gateway command reset the entry point for container:

```text
docker run -v $PWD/data:/data epik-protocol/gateway:v0.1.0 --entrypoint=gateway version
```

