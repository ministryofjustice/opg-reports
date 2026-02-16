
## Structure


## Creating a data set

Creating a complete new data set for reporting, including import and seeding of data.

**For this example we'll use `releases` as the package name**

### Create a package folder

### Package import model

- create folder & file
- naming convention of fields / json

### Database migration sql

- create const with sql
- add to _MIGRATIONS
- add tests

### Import sql

- create const
- add to _IMPORTS
- update statementForT switch
- add to tests

### Seeds

- create generation function, returns set of import model
- add to SeedAll
- add to tests in package
- add to tests in seed cmd


### Create data getter

- within package sub-folder
- returns slice of model
- add new clients where needed

### Create an import sub-command

- folder in package
- copy similar one and adjust
- add to the main import command

### API folder

- make `{x}dynamic` folder & empty file
- copy from similar endpoint
- update endpoint details at start
- update query / request / filter / response
- create a new model for the data
- update the query segments (sql builder) blocks to match needs



## Running imports

```
./import teams --db="../api/database/api.db"
./import accounts --db="../api/database/api.db"
./import codebases --db="../api/database/api.db"
./import codeowners --db="../api/database/api.db"
aws-vault exec ${profile} -- ./import infracosts --db="../api/database/api.db"
aws-vault exec ${profile} -- ./import uptime --db="../api/database/api.db"
```

----

SELECT
   A.service,
   A.date as dateA,
   CAST( coalesce(SUM(A.cost), 0) as REAL) as costA,
   B.date as dateB,
   CAST( coalesce(SUM(B.cost), 0) as REAL) as costB,
   (
      CAST( coalesce(SUM(A.cost), 0) as REAL) -
         CAST( coalesce(SUM(B.cost), 0) as REAL)
   ) as diff
FROM aws_costs as A
LEFT JOIN aws_costs as B ON
   B.service = A.service AND
   strftime("%Y-%m", B.date) = "2025-11"
WHERE
   A.service != 'Tax' AND
   strftime("%Y-%m", A.date) = "2025-08"
GROUP BY
   A.service,
   strftime("%Y-%m", A.date)
HAVING
   abs(diff) > 100
ORDER BY
   abs(diff) DESC

;
