// sqlx
// - http://jmoiron.github.io/sqlx/
// - https://kadirseckin.medium.com/simplify-database-operations-in-golang-with-sqlx-bbbfed23bb1f
// - https://github.com/joncrlsn/go-examples/blob/master/sqlx-sqlite.go


// openapi spec and server
// - https://github.com/danielgtaylor/huma
//      - https://huma.rocks/tutorial/writing-tests/
//      - https://go.dev/play/p/eprCn3LmPgV - uses standard router, line 68 for getting the spec

// go-faker for faking data







TODOs
    - big import change
        - create a units & units->github teams exported file
            - uploads to bucket
        - create old -> new single script to convert all old records to new structure
            - N costs
            - N uptime
        - import will then happen in api init
            - create a top level `info` folder to take over from `pkg/bi`
            - creates a db with the sqlite adaptor
            - creates all the tables with bootstrap
                - tables selected depends on `info` Mode
                - table map to a bucket name & ordered model list
            - if `info` Dataset is set to fake
                - create fake data and write to files in folder
            - if `info` Dataset is set to real
                - download bucket data to local folder
            - import then reads content of each folder
    - add info logging when registering to huma in each source .Register func
    - testing on the collectors



