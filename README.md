# opg-reports

This repository acts as a central hub to generate, combine and display series of key data that we report on frequently to both internal and external parties.


## Running Development Environment

Using `go` and binaries directly you can run dev environment directly.

To generate a fresh build (in `./builds`) run:

```bash
make clean && make go-build
```

This will remove any existing built binaries and create new built versions as well as setting up folder structures and assets.

The `goreleaser` process will copy assets for the front and api servers into correct locations and try to download real data from the s3 bucket.

To then run the api:

```bash
make go-run-api
```

On startup, the api will look for databases, if those aren't found but it has csv and schema files then it will create a database from those. Otherwise, it will generate a randomly seeded set of databases.

Then run the front end with:

```bash
make go-run-front
```

The front end will fetch govuk-frontend assets as part of its startup.
