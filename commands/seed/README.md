# Seed

The seed command handles importing / generating databases for the [`api`](../../servers/api/). If the data file pattern matches and finds files, then those are imported; otherwise it will try and generate dummy data.

## Example

```bash
go run main.go \
    -table github_standards \
    -db ./builds/api/github_standards.db \
    -schema ./builds/api/github_standards/github_standards.sql \
    -data "./builds/api/github_standards/data/*.json"
```

- `table`: the name of the table to import data into
  - *Note:* this is also used as the key to determine how to generate / insert data, so has to be a known value
- `-db`: the path to create the database at. If this exists, the command will exit
- `schema`: the database schema sql file to use to create the database
- `data`: a file pattern to find data files. If these don't exist, fake table will be created

