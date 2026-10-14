# Integra Core (Open-Core SDK)

[![Go Reference](https://pkg.go.dev/badge/github.com/mateusz-kacpura/integra-core.svg)](https://pkg.go.dev/github.com/mateusz-kacpura/integra-core)
[![Documentation Portal](https://img.shields.io/badge/docs-portal-3b82f6?style=flat-square&logo=gitbook&logoColor=white)](https://mateusz-kacpura.github.io/integra-core/)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

🌐 **[Live Documentation Portal](https://mateusz-kacpura.github.io/integra-core/)** – Interaktywny portal dokumentacji architektonicznej, wersja dwujęzyczna (PL/EN), pełnoekranowe schematy Mermaid i wyszukiwarka na żywo.

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
