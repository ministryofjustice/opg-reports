--
CREATE TABLE aws_uptime (
    -- auto inc primary key
    id INTEGER PRIMARY KEY,
    -- timestamp for when we generated this record
    ts TEXT NOT NULL,
    --
    unit TEXT NOT NULL,
    average REAL NOT NULL,
    date TEXT NOT NULL
    --
) STRICT;

-- used to track the dates of when the data was pulled
CREATE TABLE aws_uptime_tracker (
    id INTEGER PRIMARY KEY,
    run_date TEXT NOT NULL
) STRICT ;

--
CREATE INDEX aws_uptime_date_idx ON aws_costs(date);
