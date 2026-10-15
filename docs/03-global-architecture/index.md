# 3. Global System Architecture and Mermaid Visualizations

[◀ Previous: Benchmark Analysis](../02-benchmark-analysis-and-performance/index.md) | [▲ Master Table of Contents](../index.md) | [Next: Key System Workflows ▶](../04-key-system-workflows/index.md)

---

> [!TIP]
> **Interactive Fullscreen Diagrams**: You can also open the standalone [DIAGRAMS-VIEWER.html](../DIAGRAMS-VIEWER.html) to interactively zoom, pan, and view all system architecture diagrams in full screen.

## 📑 Subsections in this Chapter

- [**3.1. Ecosystem and Infrastructure Architecture**](3.1-ecosystem-and-infrastructure.md)
- [**3.2. Go Workspace Dependencies and Open-Core Split**](3.2-go-workspace-and-open-core.md)
- [**3.3. Internal Clean Architecture + CQRS of Microservices**](3.3-clean-architecture-and-cqrs.md)

---

# 3.1. Ecosystem and Infrastructure Architecture (Edge, Gateways, Services, Data & Observability)

[◀ Back to Chapter](index.md) | [▲ Master Table of Contents](../index.md)

---

The diagram below illustrates the full production deployment topology defined in `integra-infra/` (encompassing the **BunkerWeb + ModSecurity** WAF edge cluster, three isolated **Envoy/Gateway** entrypoints, **Linkerd Service Mesh**, Go microservices, database clusters, and telemetry):

```mermaid
graph TB
    subgraph External["External Clients & Devices"]
        WebClient["Web / Mobile Apps<br/>(api.integra.com)"]
        AdminClient["Admin Portal<br/>(admin.integra.com)"]
        IoTDevice["NFC / IoT Terminals<br/>(iot.integra.com)"]
    end

    subgraph EdgeLayer["Edge Infrastructure (integra-infra/deploy/edge)"]
        BunkerWeb["BunkerWeb WAF + ModSecurity<br/>(AntiBot, Rate Limit 100r/s, IP Blacklist, Let's Encrypt TLS/HTTP2)"]
        PubGW["Public Gateway<br/>(:8080)"]
        AdmGW["Admin Gateway<br/>(:8081)"]
        IoTGW["IoT Gateway<br/>(:8082)"]
    end

    subgraph K8sMesh["Kubernetes / EKS Cluster + Linkerd mTLS Service Mesh (namespace: integra)"]
        subgraph IAM["integra-iam-service"]
            IAM_Srv["IAM Server<br/>(HTTP :8081 / gRPC :9081)"]
            IAM_Wrk["IAM Outbox & Purge Worker"]
            IAM_Adm["IAM Admin CLI"]
        end

        subgraph USER["integra-user-service"]
            USR_Srv["User Server<br/>(HTTP :8082 / gRPC :9082)"]
            USR_Wrk["User Kafka & Outbox Worker"]
        end

        subgraph LEDGER["integra-ledger-service"]
            LDG_Srv["Ledger Server<br/>(HTTP :8085 / gRPC :9085)"]
            LDG_Wrk["Ledger Outbox & Snapshot Worker"]
            LDG_Cat["Project-Catchup CLI"]
        end
    end

    subgraph DataLayer["Data & Messaging Tier (AWS RDS, DocumentDB, ElastiCache, MSK)"]
        PG_IAM[("PostgreSQL 18<br/>Schema: iam")]
        PG_USR[("PostgreSQL 18<br/>Schema: users")]
        PG_LDG[("PostgreSQL 18<br/>Schema: ledger")]
        Mongo[("MongoDB 8 / DocumentDB<br/>(Auth Sessions & Password Reset)")]
        Valkey[("Valkey 9 / ElastiCache<br/>(L2 Security State, Leases, Claims Cache)")]
        Kafka{{"Apache Kafka / AWS MSK<br/>(38 Canonical Event Streams)"}}
    end

    subgraph ObsLayer["Observability Tier (integra-infra/monitoring)"]
        OTel["Jaeger / OpenTelemetry<br/>(W3C Distributed Tracing :4318)"]
        Prom["Prometheus + Alertmanager<br/>(Metrics & SLA Alert Rules)"]
        Loki["Fluentd + Grafana Loki<br/>(PII-Sanitized Logs)"]
        Grafana["Grafana Dashboards<br/>(Overview, IAM, User, Ledger)"]
    end

    WebClient --> BunkerWeb
    AdminClient --> BunkerWeb
    IoTDevice --> BunkerWeb

    BunkerWeb --> PubGW
    BunkerWeb --> AdmGW
    BunkerWeb --> IoTGW

    PubGW -. "mTLS" .-> USR_Srv
    PubGW -. "mTLS" .-> IAM_Srv
    PubGW -. "mTLS" .-> LDG_Srv

    AdmGW -. "mTLS Admin" .-> IAM_Srv
    AdmGW -. "mTLS Admin" .-> LDG_Srv
    IoTGW -. "mTLS IoT" .-> LDG_Srv

    IAM_Srv --> PG_IAM
    IAM_Srv --> Mongo
    IAM_Srv --> Valkey
    IAM_Wrk --> PG_IAM
    IAM_Wrk --> Kafka

    USR_Srv --> PG_USR
    USR_Srv --> Valkey
    USR_Wrk --> Kafka
    USR_Wrk --> PG_USR

    LDG_Srv --> PG_LDG
    LDG_Srv --> Valkey
    LDG_Wrk --> PG_LDG
    LDG_Wrk --> Kafka

    K8sMesh -. "W3C Traces" .-> OTel
    K8sMesh -. "Prometheus Metrics" .-> Prom
    K8sMesh -. "Sanitized JSON Logs" .-> Loki

    Prom --> Grafana
    Loki --> Grafana
    OTel --> Grafana
```

---

[◀ Back to Chapter](index.md) | [▲ Master Table of Contents](../index.md)


# 3.2. Go Workspace (`go.work`) Dependency Architecture and Open-Core Split

[◀ Back to Chapter](index.md) | [▲ Master Table of Contents](../index.md)

---

The diagram below details the module hierarchy inside the Go 1.26 Workspace (`go.work`) and the automated extraction barrier defining the Open-Core layer:

```mermaid
graph TD
    subgraph Monorepo["Go 1.26 Workspace (go.work)"]
        Examples["examples<br/>(Executable Demos)"]
        IAM["integra-iam-service<br/>(Auth, RBAC, Licenses)"]
        User["integra-user-service<br/>(Profiles & Skills)"]
        Ledger["integra-ledger-service<br/>(Double-Entry Accounting)"]
        Shared["integra-shared-go<br/>(Enterprise Shared Core)"]
        Adapters["integra-postgres-adapters<br/>(SQL Routing & Multi-Tenancy)"]
        Proto["integra-proto<br/>(Protobuf / ConnectRPC)"]
    end

    subgraph OpenCoreDist["Extracted Open-Core Distribution (.dist/integra-core)"]
        CoreDist["github.com/mateusz-kacpura/integra-core<br/>(Apache 2.0 Open-Core Module)"]
    end

    Examples --> Shared
    Examples --> Proto
    Examples --> Adapters

    IAM --> Shared
    IAM --> Proto
    IAM --> Adapters

    User --> Shared
    User --> Proto
    User --> Adapters

    Ledger --> Shared
    Ledger --> Proto
    Ledger --> Adapters

    Adapters --> Shared

    Shared -. "Extract via Prepare-OpenCoreSplit.py" .-> CoreDist
    Proto -. "Extract via Prepare-OpenCoreSplit.py" .-> CoreDist
```

---

[◀ Back to Chapter](index.md) | [▲ Master Table of Contents](../index.md)


# 3.3. Internal Clean / Hexagonal Architecture + CQRS of Microservices

[◀ Back to Chapter](index.md) | [▲ Master Table of Contents](../index.md)

---

Every microservice (`integra-iam-service`, `integra-user-service`, `integra-ledger-service`) enforces a strict Hexagonal / Clean Architecture model with CQRS command and query separation:

```mermaid
flowchart TB
    subgraph Transport["Transport Layer (internal/transport)"]
        REST["HTTP / REST Handlers<br/>(Gin Framework)"]
        gRPC["gRPC / ConnectRPC<br/>(Protobuf Handlers)"]
        KafkaConsumer["Kafka Event Consumers<br/>(Sarama / Segmentio)"]
    end

    subgraph Application["Application Layer (internal/application)"]
        Commands["Commands / Write Path (CQRS)<br/>(RegisterUser, CommitPosting, CreateInvoice)"]
        Queries["Queries / Read Path (CQRS)<br/>(GetProfile, ListLedgerEntries, VerifyToken)"]
        Ports["Outbound Ports (Interfaces)<br/>(UserRepository, LedgerRepository, SecurityState)"]
    end

    subgraph Domain["Domain Layer (internal/domain)"]
        Entities["Entities & Value Objects<br/>(User, Money, Posting, Account)"]
        Aggregates["Aggregates (Event Sourcing)<br/>(LedgerAggregate, TenantLicense)"]
        DomainEvents["Domain Events<br/>(AccountCredited, UserRegistered)"]
    end

    subgraph Infrastructure["Infrastructure Layer (internal/infrastructure)"]
        PG["PostgreSQL Adapters<br/>(pgx/v5, Single CTE SQL)"]
        VK["Valkey L2 Leases<br/>(Lua Scripts, Lease Checks)"]
        OB["Outbox Publisher<br/>(Transactional Outbox Table)"]
    end

    Transport --> Application
    Application --> Domain
    Infrastructure -. "implements" .-> Ports
    VK -. "implements" .-> Ports
    OB --> PG
```

---

[◀ Back to Chapter](index.md) | [▲ Master Table of Contents](../index.md)


---

[◀ Previous: Benchmark Analysis](../02-benchmark-analysis-and-performance/index.md) | [▲ Master Table of Contents](../index.md) | [Next: Key System Workflows ▶](../04-key-system-workflows/index.md)
