# distributed-log-engine

**How to build it in phases**

Phase 1 — Single-node write-ahead log. A simple Go server that accepts log entries over HTTP or TCP, appends them to a binary file with offsets, and lets consumers read from any offset. No partitions yet, no replication. Just the log. Get this working and shippable.

Phase 2 — Topics and partitions. Introduce the concept of named topics. Each topic gets N partition files. A producer hashes a key to pick a partition. Consumers track their offset per partition. This is where Go's goroutine model starts to shine — one goroutine per partition reader.

Phase 3 — Consumer groups. Multiple consumers can join a group and the partitions are distributed between them. If one drops, partitions rebalance. This is the hard, interesting part — and it's the thing every platform engineer is expected to understand.

Phase 4 (optional, bring in Rust) — The serialization and parsing hot path. Message encoding (protobuf or MessagePack), checksum validation, and binary format parsing are a clean boundary where Rust earns its place without you rewriting the whole project. A Rust crate that your Go server calls via FFI or a sidecar process.