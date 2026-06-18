# distributed-log-engine

A distributed write-ahead log built in Go, used as the event backbone for a self-hosted Strava activity heatmap.

Every time I complete an activity, Strava fires a webhook. That webhook is a producer — it appends an event to the log. A consumer reads from the log and updates a self-hosted heatmap rendered on OpenStreetMap tiles. Replaying the log from offset 0 rebuilds the heatmap from scratch. That's the whole idea.

```
Strava activity
      │
      ▼
Strava Webhook ──► Log Engine (producer) ──► logs.bin / logs.idx
                                                      │
                                          ┌───────────┴───────────┐
                                          ▼                       ▼
                                   Heatmap Consumer       Stats Consumer
                                  (MapLibre GL map)    (mileage, pace, PRs)
```

---

## Phases

**Phase 1 — Single-node write-ahead log**
A Go server that accepts log entries over HTTP, appends them to a binary file with offsets, and lets consumers read from any offset. The Strava webhook hits a `/append` endpoint; the entry lands in `logs.bin` with a checksum and its byte offset recorded in `logs.idx`. No partitions yet — just the log.

**Phase 2 — Topics and partitions**
Introduce named topics. Activity events go to an `activities` topic, split across N partition files. A producer hashes the activity ID to pick a partition. Consumers track their offset per partition. One goroutine per partition reader — this is where Go's concurrency model starts to earn its place.

**Phase 3 — Consumer groups**
Multiple consumers join a group and partitions are distributed between them. The heatmap consumer and the stats consumer run independently, each tracking their own offsets. If one crashes, partitions rebalance. This is the hard part — and the thing every platform engineer is expected to understand.

**Phase 4 (optional) — Rust on the hot path**
Message encoding (protobuf or MessagePack), checksum validation, and binary format parsing are a clean boundary where Rust earns its place without rewriting the whole project. A Rust crate called via FFI or a sidecar process handles serialization; Go handles everything else.

---

## Why a log and not a database?

A database stores current state. A log stores what happened and when. The heatmap is just a materialized view — one way of reading the same sequence of events. If I want a different view (weekly mileage, elevation gain, PR history), I add a new consumer and replay from offset 0. The source of truth never changes.

---

## Self-hosted stack

- **Log engine** — this repo
- **Map tiles** — OpenStreetMap via a self-hosted tile server
- **Heatmap renderer** — MapLibre GL, decoding Strava's encoded polylines per activity
- **Strava webhook** — registered via the Strava API, forwarded to the log's `/append` endpoint
