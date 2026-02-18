
## Structure

using `account`

```
mkdir -p ./report/internal/account
mkdir -p ./report/internal/account/accountapi
mkdir -p ./report/internal/account/accountimport
```

migrations:
   add db creation / migration sql to `./report/internal/global/migrations/statements.go`
   add map entry in `./report/internal/global/migrations/migrations.go`
   update `MigrateAll` func to include any new segments

getting raw data (import):
   create raw data importer - `./report/internal/account/accountimport/import.go`; copy similar data source
   update the `InsertStatement` string to be accurate
   update the `Model` struct
   update the `Args` struct if required
   adjust `Import` func
   add tests

add seeds
   add new seed generation in `./report/internal/global/seeds/seeds.go`
   update `SeedAll` to include the new generator
   add propertry to the results struct
   add check in seeds test

main import command
   in the main import command `./report/cmd/import/main.go` add a new subcommand and func
   func should call migration; then run the previous import command
   test sub command

build the api endpoints


add api enpoint register to main api cmd
