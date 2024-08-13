# opg-reports

This repository acts as a central hub to generate, combine and display series of key data that we report on frequently to both internal and external parties.


## Development environment<a name="development-environment"></a>

### Using `go` directly<a name="development-environment-go"></a>

Using `go` and binaries directly you can run dev environment directly.

To generate a fresh build (in `./builds`) run:

```bash
make clean && make go-build
```

This will remove any existing built binaries and create new built versions as well as setting up folder structures and assets.

The `goreleaser` process will copy assets for the front and api servers into correct locations and try to download real data from the s3 bucket.

To then run the api:

```bash
make go-up-api
```

On startup, the api will look for databases, if those aren't found but it has csv and schema files then it will create a database from those. Otherwise, it will generate a randomly seeded set of databases.

Then run the front end with:

```bash
make go-up-front
```

The front end will fetch govuk-frontend assets as part of its startup.

### Using `docker compose`<a name="development-environment-docker"></a>

You can spin up versions of the code base using the provided docker compose files (`docker-compose.yml` and `./docker/docker-compose.dev.yml`). We have targets in the makefile for this, to build the local images for development, run:

```bash
make clean && make docker-up
```

If you want to build the production versions, run:

```bash
make clean && make docker-up-production
```


