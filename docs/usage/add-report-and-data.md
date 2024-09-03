# Adding a report and associated data

- Decide on the name of the report, as its used in many places (`<report_name>`)
- Initial setup:
  - Setup data structures for what data to capture:
    - create a folder for this report in `./datastore` (`./datastore/<report_name>`)
    - create the `sqlc.yaml` file and configure sources - refer to an existing version
    - create a schema file for this report in `./datastore/<report_name>/<report_name>.sql`
    - create a simple version of sql queries (most likely just inserts) in `./datastore/<report_name>/queries.sql`
    - test the schema is working by running `sqlc generate` from the new folder (`./datastore/<report_name>`)
  - Create the folder structure for new command (`./commands/<report_name>`)
    - add functionality as required for the report
    - report should output its data to a relative `./data` directory - this is used in build processes and local usage
    - report data should be a json list (using marshaling)
    - create readme on how to use the command (`./commands/<report_name>/README.md`) with examples
  - Update `Makefile`
    - Add a new `sqlc generate` command to `data/sqlc` target for this report
    - Add a new bucket download command to `data/sync` target for this report (bucket subfolder should be `<report_name>`)
  - Add build for this command to the `./.github/actions/go_build/action.yml` list
  - Check in and commit these changes
- Workflow to run the command:
  - add a workflow to run report and upload it to the bucket (use an existing one as base)
  - enable pull request triggers on the workflow for testing within the pr
  - run this, or push the data up manually
- Data importing:
  - create a fake function for the new model
  - add a `models.extras.go` into the sqlc area to do insertable method (and any other issues)
  - create new generation (fake data) handler in `./commands/seed/seeder/generators.go`
  - create new insertion (real data) handler in `./commands/seed/seeder/insertors.go`
  - create new tracker handler n `./commands/seed/seeder/trackers.go`
  - add to api Dockerfile in base and build stages


- Update datastore setup
  - add simple queries (normally an insert) to queries
- Create seeder for the data

- Add api end point
  - In `./servers/api/` folder path add a new top level folder of the name
  - use an existing api as a base
  - update the api handling to features you need
  - add new response type
  - add tests
  - update sqlc queries
- Add front handler
  - update the config.json with new values
  - copy existing version into new folder name
  - create handler for the api end points
  - create templates
