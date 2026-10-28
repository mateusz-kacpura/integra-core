# 5. Detailed Directory and File Documentation

[◀ Previous: Key System Workflows](../04-key-system-workflows/index.md) | [▲ Master Table of Contents](../index.md)

---

## 📑 Subsections in this Chapter

- [**5.1. Root Monorepo, benchmarks/, examples/, and scripts/**](5.1-root-monorepo-and-tooling.md)
- [**5.2. Module .dist/integra-core (Exported Open-Core Distribution)**](5.2-open-core-distribution.md)
- [**5.3. Modules integra-proto and integra-postgres-adapters**](5.3-proto-and-postgres-adapters.md)
- [**5.4. Module integra-shared-go (Enterprise Shared Core)**](5.4-shared-go-library.md)
- [**5.5. Microservice integra-iam-service (Identity & Licenses)**](5.5-iam-service.md)
- [**5.6. Microservice integra-user-service (Profiles & Skills)**](5.6-user-service.md)
- [**5.7. Microservice integra-ledger-service (General Ledger & Invoicing)**](5.7-ledger-service.md)
- [**5.8. Module integra-infra (Cloud & Kubernetes)**](5.8-cloud-infrastructure.md)

---

# 5.1. Root Monorepo, `benchmarks/`, `examples/`, and `scripts/`

[◀ Back to Chapter](index.md) | [▲ Master Table of Contents](../index.md)

---

| File / Directory Path | Purpose & Detailed Engineering Specification |
| --- | --- |
| `go.work`, `go.work.sum` | Go 1.26 workspace definition linking 7 local modules (`examples`, `integra-iam-service`, `integra-ledger-service`, `integra-postgres-adapters`, `integra-proto`, `integra-shared-go`, `integra-user-service`) with `replace` directives. |
| `Makefile` | Central developer task automation: `fmt`, `imports`, `tidy`, `tidy-check`, `vet`, `lint`, `staticcheck`, `test` (with `-race -shuffle=on`), `test-short`, `audit`. |
| `LICENSE` | Commercial proprietary license (*All Rights Reserved*) for the Enterprise tier. |
| `.github/workflows/ci.yaml` | GitHub Actions multi-service CI pipeline running integration tests with `postgres:18-alpine` (3 isolated instances), `mongo:8`, and `valkey/valkey:9`. |
| `.gitignore`, `.dockerignore`, `.gitattributes` | Exclusion rules for build artifacts, `.pprof` profiles, `.env` secrets, and line-ending normalization (LF for Go/SQL/Proto, CRLF for PS1). |
| `benchmarks/BenchmarkCharts.py` | Zero-dependency pure Python visualization engine generating SVG vector graphics, Mermaid markdown charts, and an interactive HTML dashboard. |
| `benchmarks/BENCHMARKS-DASHBOARD.html` | Self-contained, dark-mode glassmorphism dashboard powered by Chart.js displaying L1-L3 latency tiers, alloc metrics, and historical trends. |
| `examples/go.mod`, `examples/go.sum` | Executable demonstration code module. |
| `examples/jwt-auth/main.go` | Executable demo showing Ed25519 key generation, JWT signing, and verification with L1 CPU caching (500 ms TTL). |
| `examples/sqlrouting-l1/main.go` | Executable demo showing dynamic SQL table rewriting (`outbox_events` -> `iam.outbox`) and partition compile caching in L1. |
| `scripts/Audit-Project.py` | Quality Gate orchestrator running local module validation, `staticcheck`, and race-detector tests. |
| `scripts/Check-LocalModules.py` | Automation runner executing `go` toolchain steps in each submodule with `GOWORK=off` isolation. |
| `scripts/Collect-Benchmarks.py` | Benchmark metric collector computing medians across multiple runs, logging history to `BENCHMARK-HISTORY.csv`, and enforcing the regression gate. |
| `scripts/Prepare-OpenCoreSplit.py` | Automated extraction tool generating the standalone Apache 2.0 Open-Core repository (`.dist/integra-core`). |
| `scripts/Migrate-ToEnterprisePrefix.py` | Migration utility standardizing import paths across the monorepo to `github.com/mateusz-kacpura/integra-enterprise/...`. |
| `scripts/Export-HandwrittenFiles.ps1`, `Simulate-CI.ps1` | PowerShell helpers for codebase audits, project structure dumps, and local Windows CI simulation. |

---

[◀ Back to Chapter](index.md) | [▲ Master Table of Contents](../index.md)


# 5.2. Module `.dist/integra-core` (Exported Open-Core Distribution)

[◀ Back to Chapter](index.md) | [▲ Master Table of Contents](../index.md)

---

The public open-source core (`github.com/mateusz-kacpura/integra-core`) extracted by `Prepare-OpenCoreSplit.py` provides:

* **Configuration & Contracts (`pkg/config`, `pkg/baseconfig`, `pkg/contracts`)**:
  - `BaseConfig` loaded via `viper` with `INTEGRA_` environment prefix, HTTP port validation, and `DefaultShutdownGrace = 2s`.
  - Shared DTO structures for read replicas (`UserRead`, `OrgRead` with `Version int64`) and AST architecture lint rules.
* **Health Checks (`pkg/health`)**:
  - RFC-compliant health endpoints reporting service status, database connectivity, and readiness/liveness probes.
* **JSON Serialization (`pkg/json`)**:
  - Accelerated JSON encoding via `bytedance/sonic` with automatic fallbacks for non-x86 architectures.
* **Sanitizing Logger (`pkg/logger`)**:
  - Structured zero-alloc logger integrating `zerolog` and `log/slog` with automatic regex-based PII masking (PESEL, NIP, Credit Cards, Bearer JWTs, Emails).
* **Cursor Pagination (`pkg/pagination`)**:
  - Opaque URL-safe cursor encoding and decoding with zero heap allocation (`0 B/op`).
* **Standardized API Responses (`pkg/response`)**:
  - Uniform envelope structs (`Response[T]`, `ErrorResponse`, field validation bags).
* **Distributed Telemetry (`pkg/telemetry`)**:
  - W3C Trace Context injection and extraction (`traceparent`, `X-Correlation-ID`) across HTTP, gRPC, and message brokers.
* **Canonical Topic Names (`pkg/topics`)**:
  - Constant declarations for 38 canonical Kafka domain event streams.
* **UUID Generators (`pkg/uuidgen`, `pkg/uuidx`)**:
  - Time-ordered `UUIDv7` generators optimized for B-Tree database indexing.
* **Protobuf & ConnectRPC Client Contracts (`proto/gen/go/...`)**:
  - Full Go stubs and ConnectRPC service interfaces for `IAM`, `User`, and `Accounting`.

---

[◀ Back to Chapter](index.md) | [▲ Master Table of Contents](../index.md)


# 5.3. Modules `integra-proto` and `integra-postgres-adapters`

[◀ Back to Chapter](index.md) | [▲ Master Table of Contents](../index.md)

---

#### `integra-proto/`
* **Protocol Buffers Schemas (`proto/accounting/v1`, `proto/iam/v1`, `proto/user/v1`)**:
  - Canonical gRPC service definitions for user registration, profile query, ledger transactions, and invoice management.
* **Buf Tooling (`buf.yaml`, `buf.gen.yaml`)**:
  - Linting, breaking change detection, and multi-plugin code generation (Connect-Go, Protoc-Gen-Go).
* **Generated Stubs (`gen/go/...`)**:
  - Strongly-typed protobuf structures and HTTP/2 + gRPC Connect client bindings.

#### `integra-postgres-adapters/`
* **Dynamic SQL Routing (`sqlrouting/`)**:
  - Fast-path partition table name rewriting (`sqlrouting/table.go`) with L1 cache hit in **19.1 ns**.
  - Strict SQL injection prevention using pre-compiled white-lists.
* **Connection Pooling & Migrations (`pkg/postgres`)**:
  - Tuned `pgxpool` configurations with health probes and migration harnesses.

---

[◀ Back to Chapter](index.md) | [▲ Master Table of Contents](../index.md)


# 5.4. Module `integra-shared-go` (Enterprise Shared Core)

[◀ Back to Chapter](index.md) | [▲ Master Table of Contents](../index.md)

---

Contains enterprise-grade utilities and security implementations:

* **Authenticated Encryption (`pkg/aead`, `pkg/crypto`)**:
  - `XChaCha20-Poly1305` authenticated cipher engine with automatic key derivation.
  - Zero-allocation 64KB buffers managed via `sync.Pool` achieving **0 B/op and 0 allocs/op** for large payload encryption and decryption.
* **Multi-Tier Security & Authentication (`pkg/auth`)**:
  - `JWTVerifier`: Asymmetric Ed25519 token verification with public key discovery and non-blocking `Memo` singleflight cache (eliminating cache stampedes on L1 misses).
  - `L1Cache`: Cache-line padded (`[64]byte`) in-process lock-free cache resolving tokens in **36.5 ns** with zero-alloc string-from-slice map lookups.
  - `AccountState`: Fast-invalidation lease checker synchronizing with Valkey L2 via consolidated Lua scripts.
* **Role-Based Access Control (`pkg/rbac`)**:
  - Parallel bitmask permission matrix checking permissions in **4.4 ns**.
* **Distributed Messaging & Outbox (`pkg/kafka`, `pkg/outbox`)**:
  - Pipelined worker pool with buffered channel (`DefaultDeliveryWindow = 16`), delivering a **15.75x speedup** (latency reduced from 36.55 ms to 2.30 ms).
  - Leak-free transactional outbox worker utilizing `errgroup` and reusable `time.NewTimer` instances.
* **Security State Authority (`pkg/securitystate`)**:
  - Centralized role authority lease coordination backed by Valkey Lua scripts.

---

[◀ Back to Chapter](index.md) | [▲ Master Table of Contents](../index.md)


# 5.5. Microservice `integra-iam-service` (Identity & Licenses)

[◀ Back to Chapter](index.md) | [▲ Master Table of Contents](../index.md)

---

* **Core Responsibility**: Central identity provider, authentication server, session manager, and enterprise license guard.
* **Storage Engines**:
  - PostgreSQL 18 (`schema: iam`): User credentials, role assignments, audit logs.
  - MongoDB 8: Session store and password reset tokens.
  - Valkey 9: High-throughput L2 security lease authority.
* **Key Components**:
  - `internal/application/auth`: Login, registration, token refresh, and MFA orchestration.
  - `internal/infrastructure/messaging/outbox_purge`: Automated GDPR/RODO background purge worker.
  - `cmd/iam-admin`: Administrative CLI for emergency tenant lockout and key rotation.

---

[◀ Back to Chapter](index.md) | [▲ Master Table of Contents](../index.md)


# 5.6. Microservice `integra-user-service` (Profiles & Skills)

[◀ Back to Chapter](index.md) | [▲ Master Table of Contents](../index.md)

---

* **Core Responsibility**: User profiles, account management, competency matrices, and KRS business registry sync.
* **Performance Architecture**:
  - Read replica projection caching delivering profile reads in **20.9 ns**.
  - No-Op Update Filter in `internal/app/commands/profiles.go`: Bypasses redundant database transactions and WAL writes when profile fields remain unchanged.
  - Data-Oriented Design value slices eliminating heap allocations during bulk queries.
* **Event Handlers**:
  - Consumes `canonical.iam.user_registered` to initialize default profile projections and skills vectors.

---

[◀ Back to Chapter](index.md) | [▲ Master Table of Contents](../index.md)


# 5.7. Microservice `integra-ledger-service` (General Ledger & Invoicing)

[◀ Back to Chapter](index.md) | [▲ Master Table of Contents](../index.md)

---

* **Core Responsibility**: Authoritative immutable double-entry financial accounting, real-time balance tracking, multi-currency wallets, invoices, and payouts.
* **Architectural Pillars**:
  - **TigerBeetle Authoritative Ledger Engine (`internal/infrastructure/tigerbeetle/`)**:
    - `client.go`: Direct I/O binary IPC client connecting to the TigerBeetle cluster on port 3000.
    - `id.go`: Zero-allocation bidirectional mapping between RFC 4122 UUIDs and TigerBeetle native 128-bit unsigned integers (`UUIDToUint128`, `Uint128ToUUID`).
    - `repository.go`: Implementation of `AccountRepository` and `PostingRepository` with atomic debit/credit invariants, native idempotency by Transfer ID, eliminating SQL row locks (`SELECT ... FOR UPDATE`) and database contention.
  - **Memory Arenas & Stack Hashing**: SHA-256 request frames hashed in `[512]byte` stack buffers (`0 B/op`).
  - **Pipelined Outbox Relay**: Transactional outbox worker with `DeliveryWindow = 16` publishing canonical ledger events to Kafka.
  - `cmd/project-catchup`: Production projection rebuild tool.

---

[◀ Back to Chapter](index.md) | [▲ Master Table of Contents](../index.md)


# 5.8. Module `integra-infra` (Cloud & Kubernetes)

[◀ Back to Chapter](index.md) | [▲ Master Table of Contents](../index.md)

---

* **Edge Cluster (`deploy/edge`)**: BunkerWeb WAF, ModSecurity CRS, automatic Let's Encrypt TLS renewal, rate limiting (100 req/s), and antibot rules.
* **Kubernetes & Helm (`deploy/k8s`, `deploy/helm`)**: Full Helm charts for production microservices, Linkerd mTLS service mesh, TigerBeetle StatefulSet (`deploy/k8s/infra/tigerbeetle.yaml` with Direct I/O and `IPC_LOCK`), and horizontal pod autoscalers (HPA).
* **Terraform IaC (`deploy/terraform`)**: AWS infrastructure provisioning (EKS, RDS PostgreSQL 18, ElastiCache Valkey, MSK Kafka, DocumentDB).
* **Observability Stack (`monitoring/`)**: Jaeger OpenTelemetry distributed tracing, Prometheus alert rules, Grafana Loki log aggregation, and customized dashboards.

---

[◀ Back to Chapter](index.md) | [▲ Master Table of Contents](../index.md)


---

[◀ Previous: Key System Workflows](../04-key-system-workflows/index.md) | [▲ Master Table of Contents](../index.md)
