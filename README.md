# Integra Core (Open-Core SDK)

`integra-core` is the foundational open-source SDK and contract library for the **Integra Platform**.

It provides core interface contracts, ConnectRPC/Protobuf protocols, structured logging with PII sanitization, W3C trace/correlation context propagation, cursor pagination, and essential runtime utilities.

## Installation

```bash
go get github.com/mateusz-kacpura/integra-core@v0.1.0
```

## Structure

- `proto/`: Generated ConnectRPC / Protobuf v1 contracts (`accounting`, `iam`, `user`).
- `pkg/contracts/`: Data, audit, and repository abstraction interfaces.
- `pkg/logger/`: Structured logging with PII scrubbing and correlation propagation.
- `pkg/telemetry/`: W3C distributed trace context and correlation ID extraction/injection.
- `pkg/response/`: Standardized REST / JSON HTTP responses.
- `pkg/pagination/`: Cursor-based pagination utilities.
- `pkg/uuidgen/`: UUID v7 and v4 generators.
- `pkg/health/`: HTTP liveness and readiness health checks.
- `pkg/json/`: High-performance JSON codecs.
- `pkg/topics/`: Canonical Kafka topic taxonomy.

## License

Apache License 2.0. See [LICENSE](LICENSE) for details.
