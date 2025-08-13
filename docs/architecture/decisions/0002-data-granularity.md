# 2. Data Granularity

Date: 2025-08-13

## Status

Accepted

## Context

Most data was being recorded per day, with the data importing of that data being run with varying frequency.

This led to very large database tables (millions of entries) which caused significant performance issues for running the api - as the sqlite db is stored within.

## Decision

We will now be storing data at a monthly level, but that data is updated daily, with the entire month being requested each time.

Where applicable, we'll show a notification that the data might change if it is not stable (costs etc)

## Consequences

Database insert must become upserts with valid on conflict handling.

Front end interface must indicate any areas of instability in the data shown.

