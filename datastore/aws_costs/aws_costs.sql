--
CREATE TABLE aws_costs (
    -- auto inc primary key
    id INTEGER PRIMARY KEY,
    -- timestamp for when we generated this record
    ts TEXT NOT NULL,
    --
    organisation TEXT NOT NULL,
    account_id TEXT NOT NULL,
    account_name TEXT NOT NULL,
    unit TEXT NOT NULL,
    label TEXT NOT NULL,
    environment TEXT NOT NULL,
    --
    service TEXT NOT NULL,
    region TEXT NOT NULL,
    date TEXT NOT NULL,
    cost TEXT NOT NULL,
    FOREIGN KEY(account_key) REFERENCES aws_account(id)

) STRICT;

-- used to track the dates of when the data was pulled
CREATE TABLE aws_costs_tracker (
    id INTEGER PRIMARY KEY,
    run_date TEXT NOT NULL
) STRICT ;
