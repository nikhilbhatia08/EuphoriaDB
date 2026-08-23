<p align="center">
  <img src="euphoriadb-logo.svg" alt="EuphoriaDB logo" width="480">
</p>

<h1 align="center">EuphoriaDB</h1>

<p align="center">
  A relational database engine built from scratch in Go, inspired by the classic SimpleDB architecture.
</p>

<p align="center">
  <a href="https://github.com/nikhilbhatia08/EuphoriaDB/actions"><img src="https://img.shields.io/badge/build-passing-brightgreen" alt="build status"></a>
  <img src="https://img.shields.io/badge/go-1.2x-00ADD8?logo=go&logoColor=white" alt="go version">
  <img src="https://img.shields.io/badge/status-in%20progress-yellow" alt="status">
</p>

---

## About

EuphoriaDB is a hand-built relational database management system written in Go. It follows the layered design popularized by *Database Design and Implementation* (Edward Sciore's SimpleDB), where each package builds on the one below it — from raw disk blocks all the way up to a SQL query planner.

The project is a work in progress, built as a deep dive into how databases actually work under the hood: file management, logging, buffering, transactions, record layout, query processing, parsing, and planning.

## Architecture

EuphoriaDB is organized as a stack of packages, each with a single responsibility:

| Layer | Package | Responsibility |
|---|---|---|
| Storage | `filemgr` | Manages OS files as a virtual disk of fixed-size blocks |
| Storage | `log` | Manages the write-ahead log |
| Storage | `buffer` | Manages an in-memory buffer pool that caches disk blocks |
| Transactions | `transactions` | Implements transactions at the page level (locking + logging) |
| Records | `record` | Implements fixed-length records inside pages |
| Metadata | `metadata` | Maintains metadata in the system catalog |
| Query | `query` | Implements relational algebra operators as composable scans |
| Query | `scan` | Scan interfaces used to traverse query results |
| Query | `types` | Shared value/type definitions |
| Table | `table` | Table-level abstractions over records |
| SQL | `parse` | Parses SQL statements |
| SQL | `plan` / `plan_impl` | A query planner that turns parsed SQL into a scan tree |
| Access | `driver` | Client-facing driver for connecting to the database |
| Access | `server` | Server startup and initialization |

Each layer only depends on the layers beneath it, which keeps the codebase easy to reason about and makes it a useful reference for learning how relational databases are implemented.

### Planned / in-progress

The basic engine above is intentionally simple and not optimized. The following are planned to make query processing more efficient:

- **index** — static hash and B-tree indexes, plus parser/planner support for using them
- **materialize** — materialize, sort, group-by, and merge-join operators
- **multibuffer** — multi-buffer sort and product operators for better buffer utilization
- **opt** — a heuristic-based query optimizer

## Getting started

### Prerequisites

- Go 1.2x or later ([install Go](https://go.dev/doc/install))

### Clone the repository

```bash
git clone https://github.com/nikhilbhatia08/EuphoriaDB.git
cd EuphoriaDB
```

### Build

```bash
go build ./...
```

### Try an example

```bash
go run ./examples/example1
```

## Project status

EuphoriaDB is under active development. Core storage, transaction, record, and query layers are in place; indexing, optimized operators, and the cost-based optimizer are still being built out. Expect breaking changes as the project evolves.

## Contributing

Issues and pull requests are welcome. If you're interested in database internals and want to help build out indexing, optimization, or the SQL layer, feel free to open an issue to discuss.

## License

No license has been specified yet for this repository. Until one is added, please contact the repository owner before reusing this code.

## Acknowledgements

The overall architecture is heavily inspired by the SimpleDB design described in *Database Design and Implementation* by Edward Sciore.