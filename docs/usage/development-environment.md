# Development environment

Default endpoints:

- API: [`http://localhost:8081`](http://localhost:8081)
- FRONT: [`http://localhost:8080`](http://localhost:8080)

## Using `docker`

You can spin up versions of the code base using the provided docker compose files (`docker-compose.yml` and `./docker/docker-compose.dev.yml`). We have targets in the makefile for this, to build the local images for development, run:

```bash
make clean && make docker-up
```

