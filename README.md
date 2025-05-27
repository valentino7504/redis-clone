# Go Redis Clone

A simple Redis server clone implemented in Golang.  
This project supports a subset of Redis commands and uses the RESP protocol for communication.

## Features

- RESP protocol parser
- Generic Commands
  - `PING`
  - `DEL`
  - `EXISTS`
  - `TYPE`
- String Commands
  - `SET`
  - `GET`
- Hash Commands
  - `HSET`
  - `HGET`
  - `HGETALL`
  - `HDEL`
- Int Commands
  - `INCR`
  - `DECR`
  - `INCRBY`
  - `DECRBY`
- AOF (Append Only File) persistence

## File Structure

- `main.go` — Entry point of the server
- `generic_handlers.go` — Command handlers for supported Redis generic commands
- `hash_handlers.go` — Command handlers for supported Redis hash commands
- `string_handlers.go` — Command handlers for supported Redis string commands
- `aof.go` — Handles AOF-based data persistence
- `resp.go` — RESP protocol parser and encoder
- `error_messages.go` — Error handler

## Getting Started

```bash
git clone https://github.com/valentino7504/redis-clone.git
cd redis-clone
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
