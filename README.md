# queryprofiler
Compare query profiles of 2 different servers by querying performance_schema.events_statements_summary_by_digest
# overview
The intention here is to connect to two servers and collect information from the digest table
in parallel.
This will do the following:
* generate n collections of query digests for each server.
* From this data we can collect n-1 samples which are based on collection x compared against collection x-1.
* Finally we find the top Z queries for server1
* With each query on server1 attempt to compare each of these queries with server2, comparing metrics using the available samples.

Note: this is still work in progress and not completed.
# Usage

```
queryprofiler [<options>] <dsn1> [<dsn2> ...]
```

```
DSN1='user:password@tcp(server1.example.com:3306)/performance_schema'
DSN2='user:password@tcp(server2.example.com:3306)/performance_schema'
./queryprofiler "$DSN1" "$DSN2"
```

# Concepts

In theory using P_S to profile the queries may seem quite simple,
but I think that to get useful values it requires a little more attention.
The sections below describe how queryprofiler analyses the queries
on the server.

## Event

Event is the table represetation of the P_S digest table.

## Collection

Collection is a slice of Events together with a timestamp of when
data was collected.  Collections is a slice of Collection.

## Sample

Sample is a slice of rows that come from subtracting _matching_
values by Key and recording the start time and duration of the
sample. It contains several rows for different queries.  Samples
is a slice of Sample.  Sample metrics are normalised to metrics
per second for consistency.

## Metric

This is a slice of float64, which is the underlying numbers used
by this program. Thus a sample really contains a named set of
Mmetric.

## Key

In theory the QUERY_DIGEST might be used but this digest is not
stable between different MySQL versions so I collect an MD5 digest
of the DIGEST_TEXT.  That said the DIGEST is not a unique key,
what's unique is a combination of query (digest) and SCHEMA_NAME,
so the Key considered as the key of queries is based on the MD5_DIGEST
and the SCHEMA_NAME, joined by a ".". if SCHEMA_NAME contains a
value.

## Issues

* events_statements_summary_by_digest may have empty DIGEST/DIGEST_TEXT. This represents lost values because the maximum number of digest values has been exceeded. You may see this empty query having quite high values because of this.

* Only completed queries are shown. Any long query that is running while queryprofiler is looking for data won't be shown.

* events_statements_summary_by_digest should have only one row per DIGEST_TEXT / SCHEMA_NAME. Unfortunately I've seen that this is not the case and multiple row may be present. This has been reported. See http://bugs.mysql.com/bug.php?id=79533. In the meantime if multiple rows are found with the same DIGEST_TEXT/SCHEMA_NAME the values are merged together.

* events_statements_summary_by_digest has a DIGEST column which represents a unique key (with the SCHEMA_NAME) to identify queries. However, this digest may not be the same for the same query on 2 different servers due to the way the optimiser works. Consequently queryprofiler takes an MD5 checksum of the QUERY_TEXT and uses that instead. I should really file a feature requesting that the generated query digest is calculated consistently as that would avoid this extra operation.
