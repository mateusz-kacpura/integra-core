# 3. Globalna architektura systemu i wizualizacje Mermaid

[◀ Poprzedni: Analiza benchmarków i uzasadnienie wydajnościowe](../02-benchmark-analysis-and-performance/index.md) | [▲ Główny spis treści](../index.md) | [Następny: Kluczowe przepływy pracy ▶](../04-key-system-workflows/index.md)

---

## 📑 Podrozdziały w tym dziale

- [**3.1. Architektura ekosystemu i infrastruktury**](3.1-ecosystem-and-infrastructure.md)
- [**3.2. Architektura zależności modułów Go Workspace (go.work) i podział Open-Core**](3.2-go-workspace-and-open-core.md)
- [**3.3. Wewnętrzna architektura czysta / heksagonalna + CQRS mikroserwisów**](3.3-clean-architecture-and-cqrs.md)

---

> [!TIP]
> **Interaktywny podgląd na pełnym ekranie**: Dostępna jest dedykowana przeglądarka [DIAGRAMS-VIEWER.html](../DIAGRAMS-VIEWER.html) z obsługą trybu pełnoekranowego (przycisk *Pełny ekran* lub klawisz `F11`), płynnego przybliżania kółkiem myszy (zoom) oraz przesuwania (drag & pan) wszystkich 9 schematów.

### 3.1. Architektura ekosystemu i infrastruktury (Edge, Gateways, Services, Data & Observability)

Poniższy diagram przedstawia pełną topologię wdrożeniową zdefiniowaną w katalogu `integra-infra/` (obejmującą klaster brzegowy WAF **BunkerWeb + ModSecurity**, trzy odseparowane bramy **Envoy/Gateway**, siatkę usług **Linkerd Service Mesh**, mikroserwisy Go oraz klastry bazodanowe i telemetrię):

```mermaid
graph TB
    subgraph External["Klienci zewnętrzni i urządzenia"]
        WebClient["Aplikacje Web / Mobile<br/>(api.integra.com)"]
        AdminClient["Panel Administracyjny<br/>(admin.integra.com)"]
        IoTDevice["Terminale NFC / IoT RDC<br/>(iot.integra.com)"]
    end

    subgraph EdgeLayer["Warstwa brzegowa (integra-infra/deploy/edge)"]
        BunkerWeb["BunkerWeb WAF + ModSecurity<br/>(AntiBot, Rate Limit 100r/s, IP Blacklist, Let's Encrypt TLS/HTTP2)"]
        PubGW["Public Gateway<br/>(:8080)"]
        AdmGW["Admin Gateway<br/>(:8081)"]
        IoTGW["IoT Gateway<br/>(:8082)"]
    end

    subgraph K8sMesh["Klaster Kubernetes / EKS + Linkerd mTLS Service Mesh (namespace: integra)"]
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
            LDG_Srv["Ledger Server<br/>(HTTP :8085 / gRPC :9085)<br/>+ TigerBeetle Go Client"]
            LDG_Wrk["Ledger Outbox & Pipelined Worker<br/>(DeliveryWindow = 16)"]
            LDG_Cat["Project-Catchup CLI"]
        end
    end

    subgraph DataLayer["Warstwa danych i komunikacji (TigerBeetle, AWS RDS, DocumentDB, ElastiCache, MSK)"]
        TB[("TigerBeetle Financial Ledger<br/>(Cluster ID 0, Direct I/O :3000)<br/>• debits/credits_posted<br/>• Viewstamped Replication")]
        PG_IAM[("PostgreSQL 18<br/>Schema: iam")]
        PG_USR[("PostgreSQL 18<br/>Schema: users")]
        PG_LDG[("PostgreSQL 18 (Metadane & Projekcje)<br/>Schema: ledger")]
        Mongo[("MongoDB 8 / DocumentDB<br/>(Sesje Auth & Password Reset)")]
        Valkey[("Valkey 9 / ElastiCache<br/>(L2 Security State, Leases, Claims Cache)")]
        Kafka{{"Apache Kafka / AWS MSK<br/>(38 Kanonicznych Tematów Zdarzeń)"}}
    end

    subgraph ObsLayer["Warstwa obserwowalności (integra-infra/monitoring)"]
        OTel["Jaeger / OpenTelemetry<br/>(W3C Distributed Tracing :4318)"]
        Prom["Prometheus + Alertmanager<br/>(Metryki & Reguły Alertów)"]
        Loki["Fluentd + Grafana Loki<br/>(Logi z Sanityzacją PII)"]
        Grafana["Grafana Dashboards<br/>(Overview, IAM, User, Ledger)"]
    end

    WebClient --> BunkerWeb
    AdminClient --> BunkerWeb
    IoTDevice --> BunkerWeb

    BunkerWeb --> PubGW
    BunkerWeb --> AdmGW
    BunkerWeb --> IoTGW

    PubGW --> IAM_Srv
    PubGW --> USR_Srv
    PubGW --> LDG_Srv
    AdmGW --> IAM_Srv
    IoTGW --> USR_Srv

    IAM_Srv <-->|gRPC / ConnectRPC| USR_Srv
    IAM_Srv <-->|gRPC / ConnectRPC| LDG_Srv

    IAM_Srv --> PG_IAM
    IAM_Srv --> Mongo
    IAM_Srv --> Valkey
    IAM_Wrk --> PG_IAM
    IAM_Wrk --> Kafka

    USR_Srv --> PG_USR
    USR_Srv --> Valkey
    USR_Wrk <--> Kafka
    USR_Wrk --> PG_USR

    LDG_Srv --> TB
    LDG_Srv --> PG_LDG
    LDG_Srv --> Valkey
    LDG_Wrk <--> Kafka
    LDG_Wrk --> PG_LDG
    LDG_Cat --> PG_LDG

    IAM_Srv -.-> OTel
    USR_Srv -.-> OTel
    LDG_Srv -.-> OTel
    Prom -.-> IAM_Srv
    Prom -.-> USR_Srv
    Prom -.-> LDG_Srv
    Grafana --> Prom
    Grafana --> Loki
    Grafana --> OTel

```

---

### 3.2. Architektura zależności modułów Go Workspace (`go.work`) i podział Open-Core

Wzorzec organizacji workspace'u w Go (`go.work`) pozwala na jednoczesną pracę nad wieloma powiązanymi modułami bez konieczności ich wcześniejszego publikowania w zdalnym repozytorium. Jednocześnie skrypty `scripts/Check-LocalModules.py` oraz `scripts/Audit-Project.py` weryfikują kompilację każdego modułu w izolacji z flagą `GOWORK=off`:

```mermaid
graph TD
    subgraph Monorepo["Go 1.26 Workspace (go.work)"]
        Examples["examples<br/>(jwt-auth, sqlrouting-l1)"]
        IAM["integra-iam-service"]
        USER["integra-user-service"]
        LEDGER["integra-ledger-service"]

        Shared["integra-shared-go<br/>(Auth L1/L2, Crypto ChaCha20, RBAC,<br/>Outbox, Kafka, SecurityState, Telemetry)"]
        PGAdapters["integra-postgres-adapters<br/>(sqlrouting L1 + sqlc dbgen)"]
        Proto["integra-proto<br/>(Buf v2: iam/v1, user/v1, accounting/v1)"]
    end

    subgraph OpenCore[".dist/integra-core (eksport Apache 2.0)"]
        OC_Pkg["pkg/*<br/>(config, contracts, health, jsoncodec,<br/>logger, pagination, response, telemetry, topics, uuidgen)"]
        OC_Proto["proto/*<br/>(iam/v1, user/v1, accounting/v1)"]
    end

    Examples --> Shared
    Examples --> PGAdapters

    IAM --> Shared
    IAM --> PGAdapters
    IAM --> Proto

    USER --> Shared
    USER --> PGAdapters
    USER --> Proto

    LEDGER --> Shared
    LEDGER --> PGAdapters
    LEDGER --> Proto

    Shared --> PGAdapters
    Shared --> Proto

    Shared -. "scripts/Prepare-OpenCoreSplit.py" .-> OC_Pkg
    Proto -. "scripts/Prepare-OpenCoreSplit.py" .-> OC_Proto

```

---

### 3.3. Wewnętrzna architektura czysta / heksagonalna + CQRS mikroserwisów

Każdy z trzech serwisów (`integra-iam-service`, `integra-user-service`, `integra-ledger-service`) rygorystycznie przestrzega podziału na warstwy w katalogu `internal/` (który w Go na poziomie kompilatora uniemożliwia importowanie pakietów spoza drzewa rodzica):

```mermaid
flowchart TB
    subgraph Transport["Warstwa transportu (internal/transport)"]
        REST["REST API (Gin + Sonic JSON)<br/>Middleware: Telemetry, RateLimit, Idempotency, Auth, RBAC"]
        GRPC["gRPC / ConnectRPC Server<br/>(Walidacja kontraktów Protobuf)"]
        KafkaIn["Kafka Consumer<br/>(Odbiór zdarzeń domenowych)"]
    end

    subgraph AppLayer["Warstwa aplikacji CQRS (internal/app)"]
        Access["AccessService<br/>(Autoryzacja tożsamości, tenanta i licencji)"]
        Cmds["Commands (internal/app/commands)<br/>Mutacje stanu i przypadki użycia"]
        Qrys["Queries (internal/app/queries)<br/>Odczyty projekcji i modeli"]
        Ports["Ports (internal/app/ports)<br/>Interfejsy repozytoriów, autorytetu i audytu"]
    end

    subgraph DomainLayer["Warstwa domeny (internal/domain)"]
        Agg["Agregaty i Encje<br/>(User, Profile, LedgerAggregate, Money, Invoice)"]
        Events["Zdarzenia Domenowe & Inwarianty<br/>(Zero zależności od I/O, HTTP czy SQL)"]
    end

    subgraph InfraLayer["Warstwa infrastruktury (internal/infrastructure)"]
        TB["TigerBeetle Adapter (tigerbeetle-go)<br/>AccountRepository & PostingRepository<br/>Uint128 ID Mapping, Direct I/O"]
        PG["PostgreSQL Repository (pgx/v5 + sqlc)<br/>Unit of Work / Single-CTE Atomic Commit"]
        MG["MongoDB Repository<br/>(Sesje Auth, Tokeny Refresh, Reset Haseł)"]
        VK["Valkey Adapter<br/>(L2 Claims Cache, SecurityState, Leases)"]
        OB["Transactional Outbox Worker<br/>(Publikacja at-least-once do Kafka)"]
    end

    REST --> Access
    GRPC --> Access
    KafkaIn --> Cmds
    Access --> Cmds
    Access --> Qrys
    Cmds --> Agg
    Cmds --> Events
    Cmds --> Ports
    Qrys --> Ports
    TB -. "implementuje" .-> Ports
    PG -. "implementuje" .-> Ports
    MG -. "implementuje" .-> Ports
    VK -. "implementuje" .-> Ports
    OB --> PG

```

---

[◀ Poprzedni: Analiza benchmarków i uzasadnienie wydajnościowe](../02-benchmark-analysis-and-performance/index.md) | [▲ Główny spis treści](../index.md) | [Następny: Kluczowe przepływy pracy ▶](../04-key-system-workflows/index.md)
