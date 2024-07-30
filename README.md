# opg-reports

This repository acts as a central hub to generate, combine and display series of key data that we report on frequently to both internal and external parties.

The repository has 3 core areas, data gathering, api layer and a display layer.

## Data gathering<a name="data-gathering-intro"></a>

All code responsible for generating data is based within the `./cmd/` folder path, mostly within `./cmd/report/` folder.

Currently, this repository has reports to fetch the following data:

- [Github repository standards](./cmd/report/github/standards/README.md)
- [AWS monthly costs](./cmd/report/aws/cost/monthly/README.md)

Each report is run via a [github workflow](./.github/workflows/README.md#report-workflows) using a series of arguments to fetch and then generate a single [data file](./cmd/report/README.md#filename-pattern). This data file is then uploaded to a shared s3 bucket for use by the api layer.

Each report folder contains a `Makefile` with examples of using the command.

Please see the [report readme](./cmd/report/README.md) for more details how reports are configured and utilised.

## API layer<a name="api-layer-intro"></a>

The [api layer](./services/api/README.md) creates a single webserver and delgates the handling of various endpoints to `go` code within sub-folders.

The folder paths of the handlers [should match the report-path](./cmd/report/README.md#report-path) for consistency.

The api layer contains a `./data/` directory which is where the s3 bucket is synch'd into (via a make command).

## Display layer<a name="display-layer-intro"></a>

The [display / front layer](./services/front/README.md) uses `go` webserver and its built in templating to generate html views. It has two request handlers - static and dynamic.

The static handler handle pages that are simple html / markup that don't utilise any data from the api.

The dynamic handler fetches information from the api and will then call the template stack to generate the output using teh api data.

This is configured within the [front's config.json file](./services/front/config.json).

### Configuration

The front uses its config file to set the visiable organisation name rendered in the html (default: OPG) as well as the navigation structure with where each page gets its data from.

Additionally, what repository values are checked for the baseline and extended standard checks are configured here.

Please see the [package details for more info](./services/front/cnf/cnf.go)
