# LQS (SQL backwards)

This package provides a solution to the issue of setting SQLite PRAGMA
statements for each new connection, as discussed in
[mattn/go-sqlite3 issue #1248](https://github.com/mattn/go-sqlite3/issues/1248).

## Background

SQLite drivers like `mattn/go-sqlite3` and `modernc.org/sqlite` support custom
PRAGMA statements through URI parameters. However, this approach is not
standardized across all drivers. LQS offers a driver-agnostic way to execute
PRAGMA statements for each new connection.

## Implementation

LQS uses a custom `connector` that wraps the original driver and executes
specified SQL statements (including PRAGMA statements) when a new connection is
established. This approach works with any SQL driver that implements the
`database/sql/driver` interface.

## Usage

To use LQS, import the package and use the `Open` function instead of
`sql.Open`:

```go
import (
    "github.com/jtarchie/lqs"
    _ "github.com/mattn/go-sqlite3"
)

func main() {
    db, err := lqs.Open("sqlite3", ":memory:", "PRAGMA cache_size = 1234;")
    if err != nil {
        // Handle error
    }
    defer db.Close()

    // Use db as you would normally
}
```

The `Open` function takes three parameters:

1. The driver name (e.g., "sqlite3")
2. The data source name (DSN)
3. SQL statements to be executed for each new connection (e.g., PRAGMA
   statements)

## Testing

The package includes tests to verify that PRAGMA statements are correctly
applied to new connections and that errors are properly propagated.

```bash
task
```
