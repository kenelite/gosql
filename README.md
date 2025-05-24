# gosql

A minimalist MySQL-compatible server written in Go.  
**gosql** supports basic SQL query parsing, execution, and the MySQL client/server protocol, enabling interaction with MySQL clients such as `mysql` CLI or GUI tools.

---

## Features

- MySQL client/server protocol implementation (basic)
- SQL query parsing for simple commands (`SELECT`, `INSERT`, `CREATE TABLE`, `SHOW TABLES`, etc.)
- In-memory and file-based storage backends
- Configurable via a simple `settings/server.conf` file
- Multi-client support with concurrency
- Basic result set encoding and error handling

---

## Project Structure

```
gosql/
├── cmd/
│ └── gosql/ # main.go entrypoint
├── config/ # config parsing package
├── executor/ # SQL execution engine
├── parser/ # SQL parser
├── protocol/ # MySQL protocol implementation
├── storage/ # Storage backends (memory, file)
└── settings/ # Configuration files
```




---

## Getting Started

### Prerequisites

- Go 1.20+ installed
- MySQL client or compatible GUI tool for testing

### Installation

```bash
git clone https://github.com/yourusername/gosql.git
cd gosql
go build -o gosql ./cmd/gosql
```


### Configuration
Create or modify settings/server.conf:

```bash
[server]
port = 3306
data_path = data
```
* port: TCP port where gosql listens for connections (default 3306)
* data_path: Directory path for data storage files


### Usage
Run the server:

```bash
./gosql
```



Connect with MySQL client:
```bash
mysql -h 127.0.0.1 -P 3306 -u root
```


Execute SQL queries such as:

```sql
CREATE TABLE users (id INT, name VARCHAR(64));
INSERT INTO users VALUES (1, 'Alice'), (2, 'Bob');
SELECT * FROM users;
```






## Development
* Parsing: The parser package implements a simplistic SQL parser that converts SQL text to an AST.

* Execution: The executor runs SQL commands against the storage backends.

* Storage: Supports both in-memory and file-based key-value storage.

* Protocol: Handles MySQL protocol packets, including query reading and result set writing.



## Limitations
* Only supports a subset of MySQL SQL syntax

* No authentication yet — connections are trusted

* Limited data types (mostly strings and integers)

* No transaction or concurrency control beyond basic Go concurrency

## Future Plans
- [x] Add authentication and user management
- [x] Support more SQL syntax and data types
- [x] Implement prepared statements
- [x] Improve file storage durability and indexing
- [x] Add comprehensive tests and benchmarks


## Contributing
Contributions are welcome! Feel free to submit issues or pull requests.

## License
MIT License © 2025 Kenelite

## Contact
For questions or support, please open an issue or contact kenelite.sg@gmail.com.