# Go Redis Clone

A simple Redis server clone implemented in Golang.  
This project supports a subset of Redis commands and uses the RESP protocol for communication.

## Features

- RESP protocol parser
- Basic Redis commands:
  - `PING`
  - `SET`
  - `GET`
  - `HSET`
  - `HGET`
  - `HGETALL`
- AOF (Append Only File) persistence

## File Structure

- `main.go` — Entry point of the server
- `aof.go` — Handles AOF-based data persistence
- `handler.go` — Command handlers for supported Redis commands
- `resp.go` — RESP protocol parser and encoder

## Getting Started

```bash
git clone <your-repo-url>
cd <repo-directory>
go run .
```

You can then connect using `redis-cli`:

```bash
redis-cli
```

Try running commands like:

```bash
PING
SET mykey hello
GET mykey
HSET myhash field1 value1
HGET myhash field1
HGETALL myhash
```

## Acknowledgements

Special thanks to [Ahmed Ashraf](https://www.build-redis-from-scratch.dev/en/introduction) for his amazing article, which inspired this project.
