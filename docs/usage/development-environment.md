# Development environment

Default endpoints:

- API: [`http://localhost:8081`](http://localhost:8081)
- FRONT: [`http://localhost:8080`](http://localhost:8080)

## Using `docker`

You can spin up versions of the code base using the provided docker compose files (`docker-compose.yml` and `./docker/docker-compose.dev.yml`). We have targets in the makefile for this, to build the local images for development, run:

```bash
make clean && make up
```

*Note:* Currently there is not hot reload

### Build only

You can run only the build by calling:

```bash
make build
```

You can also limit that to a set of services by adding argument:

```bash
make build SERVICES="<A> <B> <C>"
```


## Running test suite

Call the make targets for the various types of testing:

All tests:
```bash
make tests
```

Named test:
```bash
make test name="<pattern>"
```

All benchmarks:
```bash
make benchmarks
```

Named benchmark
```bash
make benchmark name="<pattern>"
```


