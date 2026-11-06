# distbuild — Build Your Own Distributed System

Code for the SystemDR YouTube series. Over 26 weeks a single-node key-value store grows into a sharded, replicated, Raft-backed platform with a Kafka-style log. Each episode is a git tag.

| Tag | Episode | What you build |
| --- | --- | --- |
| `ep01` | Build a Key-Value Database Server in Go | In-memory store, TCP server, text protocol, load test |
| `ep02` | Never Lose Data: Write-Ahead Log from Scratch | *coming next* |

## Run episode 1

```bash
go run ./cmd/kvserver -addr :7070

# in another terminal
nc localhost 7070
SET user:1 Ada Lovelace
GET user:1
DEL user:1

# load test
go run ./cmd/kvbench -clients 50 -ops 20000

# tests, including the race detector
go test -race ./...
```

## Protocol

| Command | Reply |
| --- | --- |
| `PING` | `PONG` |
| `SET <key> <value>` | `OK` |
| `GET <key>` | value or `(nil)` |
| `DEL <key>` | `1` if deleted, `0` if missing |

## Layout

```
cmd/kvserver     the server binary
cmd/kvbench      concurrent load generator
internal/store   thread-safe in-memory map (RWMutex)
internal/server  TCP listener, one goroutine per connection, command parser
```

Watch the series: https://www.youtube.com/@SystemDR
