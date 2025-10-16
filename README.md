

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

## Adding new data

Steps
- Add new capabilities to fetch required data in the relevant `repository` package
    - maybe a new repository package is required
    - create a file for the datatype within the repository
    - for AWS extend the allowed list of sdk clients via the `SupportedClients` interface
    - Add interfaces for new `RepositoryX` & `ClientX` values
- Add migration for new data type into `DB_MIGRATIONS_UP` slice in `sqlr/schema.go`
- Add new function in the `service/seed` package for the new datatype
    - create a file for datatype
    - create a private model struct & create private insert statement
    - create sample / test data slice to insert
    - create function on the `Service` struct
    - add new func into the `All` function & return data type
- Add new, basic capabilities to the `service/api` package
    - Create a new file in the api package for this data type & `${DataType}_handlers` file as well
    - in `${DataType}` file...
        - define an insert sql statement (copy from seed)
        - define a select all sql statement
    - in `${DataType}_handlers`
        - define a `${DataType}` struct with fields to handle the results of the `all` call
        - create a `GetAll${DataType}` func on the `Service[T]` struct (see others for how)
        - create an interface `${DataType}Getter` with the `GetAll${DataType}All` func & `Closer`
        - add a `Put${DataType}` func on the `Service[T]`
- Add handlers for the basic endpoints in the `cmd/api` package
    - create new folder & package for the data type
    - create a `all` file in the new package
    - create a response struct like `Get${DataType}AllResponse`
    - create a `handleGet${DataType}All[T]` func similar to others but using the correct interfaces
    - create a `RegisterGet${DataType}All[T]` func
    - add the `RegisterGet${DataType}All[T]` func to the `RegisterHandlers` func in `cmd/api/handlers.go`
- Add new importer capabilities in `cmd/importer`
    - Will use the sevice package above for inserting
    - New file for the data type
    - Create a new cobra command to handle the data type
        - have a couple of sub funcs to get the data (via repository) and then insert it (via service)
    - add a github workflow to run the command named `reports_${datatype}`
        - run daily via cron, avoid other time slots where can to reduce db clashes
- Revist and add capabilities to the `service/api` package for extended functionality
    - Add new handlers and register them for the `cmd/api` package for any new features


