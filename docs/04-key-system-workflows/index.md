# 4. Key System Workflows

[◀ Previous: Global Architecture](../03-global-architecture/index.md) | [▲ Master Table of Contents](../index.md) | [Next: Directory Structure ▶](../05-directory-and-file-structure/index.md)

---

## 📑 Subsections in this Chapter

- [**4.1. Workflow 1: Multi-Tier JWT Token and Security State Verification**](4.1-jwt-and-security-verification.md)
- [**4.2. Workflow 2: Atomic Double-Entry Ledger Posting & Idempotency**](4.2-double-entry-ledger-posting.md)
- [**4.3. Workflow 3: Distributed User Registration Saga**](4.3-user-registration-saga.md)
- [**4.4. Workflow 4: CI/CD Quality Pipeline and Regression Gate**](4.4-cicd-quality-pipeline.md)

---

# 4.1. Workflow 1: Multi-Tier JWT Token and Security State Verification (L1 -> L2 -> L3)

[◀ Back to Chapter](index.md) | [▲ Master Table of Contents](../index.md)

---

```mermaid
sequenceDiagram
    autonumber
    actor Client as API Client
    participant MW as Gin Auth Middleware
    participant L1 as L1 Cache (CPU In-Process, 500ms)
    participant VK as Valkey L2 (Lease Check, 1000ms)
    participant PG as PostgreSQL 18 L3 (User/IAM State)

    Client->>MW: HTTP Request with Bearer JWT
    MW->>MW: Verify Ed25519 Signature & Expiration

    alt Signature or Format Invalid
        MW-->>Client: 401 Unauthorized (Invalid Token)
    end

    MW->>L1: Check Security State (Account Status, Permissions)

    alt L1 Cache Hit (< 36.5 ns)
        L1-->>MW: State Valid (Cached)
        MW-->>Client: 200 OK (Proceed to Handler)
    else L1 Cache Miss (> 500 ms elapsed)
        MW->>VK: Check L2 Lease / Revocation List (Eval Lua)
        alt L2 Valkey Hit
            VK-->>MW: Lease Active & Roles Valid
            MW->>L1: Store State in L1 (500ms TTL)
            MW-->>Client: 200 OK (Proceed to Handler)
        else L2 Valkey Miss / Lease Expired
            MW->>PG: Fallback Query to L3 (Account State & Version)
            PG-->>MW: Return Database State
            MW->>VK: Populate L2 Lease (1000ms TTL)
            MW->>L1: Populate L1 Cache (500ms TTL)
            MW-->>Client: 200 OK (Proceed to Handler)
        end
    end
```

---

[◀ Back to Chapter](index.md) | [▲ Master Table of Contents](../index.md)


# 4.2. Workflow 2: Atomic Double-Entry Ledger Posting & Idempotency

[◀ Back to Chapter](index.md) | [▲ Master Table of Contents](../index.md)

---

```mermaid
sequenceDiagram
    autonumber
    actor Caller as Client / Calling Service
    participant Handler as Ledger Handler (HTTP/gRPC)
    participant Domain as Domain (Money & Uint128 Mapping)
    participant TB as TigerBeetle Cluster (Direct I/O :3000)
    participant DB as PostgreSQL 18 (Metadata & Outbox)
    participant Outbox as Pipelined Worker Pool (DeliveryWindow = 16)
    participant Kafka as Apache Kafka (finance.ledger.events)

    Caller->>Handler: PostAccounts(DebitUUID, CreditUUID, Amount, RequestID)
    Handler->>Domain: Execute In-Memory Validation (Currency matching, Minor units)
    Handler->>Domain: Map UUIDs to TigerBeetle Uint128 (UUIDToUint128, 0 B/op)
    Domain-->>Handler: Uint128 (DebitID, CreditID, TransferID)

    Handler->>TB: CreateTransfers([]types.Transfer)
    Note over TB: Deterministic Storage Engine (Direct I/O, IPC_LOCK):<br/>1. Atomic balance verification & limits check<br/>2. Native idempotency by Transfer ID<br/>3. Real-time debits/credits_posted accumulation<br/>4. Zero SQL row-locks & zero snapshot overhead
    alt Validation Failure (e.g. ExceedsCredits / Invalid Account)
        TB-->>Handler: TransferError (Result code)
        Handler-->>Caller: 422 Unprocessable Entity / 409 Conflict
    else Transfer Succeeded
        TB-->>Handler: Transfer OK (0 errors, instant confirmation)
        Handler->>DB: Store optional business metadata & outbox entry
        Handler-->>Caller: 200 OK (TransferID, updated balance)
    end

    par Asynchronous Pipelined Outbox Relay (DeliveryWindow = 16)
        Outbox->>DB: Read Unprocessed Events
        Outbox->>Kafka: Publish Canonical Event (finance.ledger.posted)
        Outbox->>DB: Mark Outbox Record as Processed
    end
```

---

[◀ Back to Chapter](index.md) | [▲ Master Table of Contents](../index.md)


# 4.3. Workflow 3: Distributed User Registration Saga and Read Replica Synchronization

[◀ Back to Chapter](index.md) | [▲ Master Table of Contents](../index.md)

---

```mermaid
flowchart LR
    subgraph IAM["1. integra-iam-service"]
        Reg["User Registration / Role Change"] --> IAM_DB[("PostgreSQL iam.users")]
        IAM_DB --> IAM_OB["Transactional Outbox"]
    end

    subgraph KafkaCluster["2. Apache Kafka Message Bus"]
        IAM_OB --> Topic{{"canonical.iam.user_registered"}}
    end

    subgraph UserSvc["3. integra-user-service"]
        Topic --> UserConsumer["Kafka Consumer"]
        UserConsumer --> UserProfile[("PostgreSQL users.profiles<br/>(Read Replica DTO)")]
        UserProfile --> UserL1["Invalidate User L1 Cache"]
    end

    subgraph LedgerSvc["4. integra-ledger-service"]
        Topic --> LedgerConsumer["Kafka Consumer"]
        LedgerConsumer --> LedgerWallet[("PostgreSQL ledger.wallets<br/>(Create Default Accounts)")]
    end
```

---

[◀ Back to Chapter](index.md) | [▲ Master Table of Contents](../index.md)


# 4.4. Workflow 4: CI/CD Quality Pipeline, Security Audit, and Benchmark Regression Gate

[◀ Back to Chapter](index.md) | [▲ Master Table of Contents](../index.md)

---

```mermaid
flowchart TD
    Push["Push / Pull Request to Monorepo"] --> Setup["CI Setup: Go 1.26, Python, Postgres 18, Valkey 9, Mongo 8, TigerBeetle 0.16"]
    Setup --> CodeCheck["Static Verification & Security:<br/>gofmt, go vet, staticcheck v0.7.0, govulncheck"]
    CodeCheck --> ModCheck["Local Module Isolation Check:<br/>scripts/Check-LocalModules.py (GOWORK=off)"]
    ModCheck --> ArchAudit["Clean Architecture & AST Audit:<br/>scripts/Audit-Project.py (go/parser)"]
    ArchAudit --> Tests["Race-Condition & Integration Tests:<br/>go test -race -tags=integration -shuffle=on ./..."]
    Tests --> Bench["Automated Benchmark Execution:<br/>scripts/Collect-Benchmarks.py"]
    Bench --> BenchGate{"Regression Barrier:<br/>Compare vs PLATFORM-BENCHMARKS-LATEST.csv<br/>Zero-Alloc Pattern Invariant Enforcement"}
    BenchGate -- "Regression <= Allowed Threshold" --> DashGen["Generate Visual Dashboard:<br/>benchmarks/BenchmarkCharts.py"]
    BenchGate -- "Regression > Threshold / Extra Allocs" --> Fail["Block Pipeline (CI Failed)"]
    DashGen --> InfraVal["Infrastructure Validation:<br/>integra-infra/scripts/validate.sh (YAML, Helm, K8s, Terraform)"]
```

---

[◀ Back to Chapter](index.md) | [▲ Master Table of Contents](../index.md)


---

[◀ Previous: Global Architecture](../03-global-architecture/index.md) | [▲ Master Table of Contents](../index.md) | [Next: Directory Structure ▶](../05-directory-and-file-structure/index.md)
