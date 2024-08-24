# Development environment

Default endpoints:

- API: [`http://localhost:8081`](http://localhost:8081)
- FRONT: [`http://localhost:8080`](http://localhost:8080)

## Running test suite

Call the make targets for the various types of testing:

All tests:
```bash
make tests
```

Named test:
```bash
make tests/named name="<pattern>"
```

All benchmarks:
```bash
make tests/benchmarks
```

Named benchmark
```bash
make tests/benchmark name="<pattern>"
```


## Using `docker`

You can spin up versions of the code base using the provided docker compose files (`docker-compose.yml` and `./docker/docker-compose.dev.yml`). We have targets in the makefile for this, to build the local images for development, run:

```bash
make clean && make docker
```

*Note:* Currently there is not hot reload

### Build only

You can run only the build by calling:

```bash
make docker/build
```

You can also limit that to a set of services by adding argument:

```bash
make docker/build SERVICES="<A> <B> <C>"
```
