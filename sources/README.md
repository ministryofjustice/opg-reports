# Sources

This series of packages represent each of the main groups of data that are being used.


## Existing sources

### costs

`costs` contains all the various database and api handlers to store and fetch all of our data relating to costs.

- `costsapi` contains all of the api handler methods and the main `Register` function to bind these to the `huma` api.
- `costsdb` containing all of the prepared database statements to be used with `datastore`.
- `costsfront` contains front end transformation handlers to convert api data into format for front end display - generally converting to table  row structure.
- `costsio` contains all of the input and output structures used by the api. Split off for ease and to avoid and circular imports.

### standards

`standards` contains everythign required to handle captgure and render of github standard data that is used. No transformer

- `standardsapi` has all of the api input and output structures as well as their
- `standardsdb` containing all the prepared statements to use with `datastore` to run queries and create database tables etc.
- `standardsio` all of the input and output structs used by the api.



