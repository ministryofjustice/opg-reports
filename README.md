

## Structure

Primary commands are with `./report/cmd/` and are the following:

- `./api/`: this command runs the api server side of the application
- `./front`: this runs the front end server for this application
- `./inmporter`: which runs commands to import data in various ways

## Design patterns

The structure of the code base revoles around service & repository pattern where the repository is responsbile for manipulation of the raw data structures and service contains the business logic.

### Repositories

The repositories are packaged based on their data source, so they each focus on a single origin, whether that's an API like GitHub or a sqlite database. They are named as `${source}r`, so `awsr` and `githubr` - where the `r` is append to avoid any naming conflicts with imported / common libraries (such as `sql`).

For reuse and mocking purposes the repositories expose and utilise a series of interfaces, which use the naming patterns of `Repository${name}` and `Client${name}`. Interfaces starting with `Repository` utilise a version of `Client` within its methods to access and manipulate the data being requested. This seperation allows either to be mocked and tested without having to make real API calls or connections to databases within the test suites.

Currently, the `sqlr` package differs and doesn't utilise `Client` interfaces, the sql connection and database are handled internally via internal functions - `connection` and `init`.

### Services

The services are packaged based on their domain area they are being used within (so `api` is used by the api commands, `front` by the front end server command) and provide functions aligned to those commands. Each service can use multiple different repository data sources within itself depending on wher ethe data needs to come from; for example the `existing` service uses both `awsr` and `githubr` to generate database entries.

The service functions exposed are aimed at solving a single, direct ask of the application ("get total cost for last month") and use both a Client and Repository to fetch that data, and then apply business logic and structural transformations within.


## Additional capabilities

There are various functional needs that are repeated within our code base that are used in multiple places; most of this code is handled under the `utils` package - this covers things like string to transformations, marshaling of structures and more.

