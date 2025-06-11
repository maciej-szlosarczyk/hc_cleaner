# HC_Cleaner

A small utility to delete stale or extra columns from honeycomb datasets.


## WARNING

This software is written just for me. Use at your own risk!

## Building

```bash
go build
./hc_cleaner
```

## Usage

You need a "Configuration Key", not an Ingest API key.

### Delete columns that have not received data in X days

```bash
hc_cleaner inactive my_dataset --since=10 --api-key=superSecret
```

### Delete columns that begin with a specific prefix

```bash
hc_cleaner prefix my_dataset --prefix=http.query_params --api-key=superSecret
```
