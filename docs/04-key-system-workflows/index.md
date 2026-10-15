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
    participant Domain as Domain Aggregate (LedgerAggregate)
    participant DB as PostgreSQL 18 (Single CTE Posting)
    participant Outbox as Transactional Outbox Worker
    participant Kafka as Apache Kafka (finance.ledger.events)

    Caller->>Handler: PostAccounts(DebitAcc, CreditAcc, Amount, RequestID)
    Handler->>Handler: Compute SHA-256 Frame Hash on Stack (0 B/op)

    Handler->>Domain: Execute In-Memory Aggregate Validation
    Note over Domain: Check debit/credit mathematical equality<br/>Check currency matching & minor unit limits
    Domain-->>Handler: Validation Successful

    Handler->>DB: Execute CommitPosting (Single CTE Transaction)
    Note over DB: 1. Verify idempotency hash (posting_requests)<br/>2. Update account balances (OCC version check)<br/>3. Insert into ledger_entries<br/>4. Append to ledger.events<br/>5. Schedule snapshot_jobs if tail >= 50<br/>6. Insert into outbox
    DB-->>Handler: Commit Confirmed (1 DB Roundtrip)

    Handler-->>Caller: 200 OK (Transaction Recorded)

    par Asynchronous Outbox Event Emission
        Outbox->>DB: Read Unprocessed Events (FOR UPDATE SKIP LOCKED)
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
    Push["Push / Pull Request to Monorepo"] --> Setup["CI Setup: Go 1.26, Python, Postgres 18, Valkey, Mongo"]
    Setup --> CodeCheck["Static Verification:<br/>gofmt, go vet, staticcheck, goimports"]
    CodeCheck --> ModCheck["Local Module Isolation Check:<br/>scripts/Check-LocalModules.py (GOWORK=off)"]
    ModCheck --> ArchAudit["Clean Architecture & AST Audit:<br/>scripts/Audit-Project.py (go/parser)"]
    ArchAudit --> Tests["Race-Condition Tests:<br/>go test -race -shuffle=on ./..."]
    Tests --> Bench["Automated Benchmark Execution:<br/>scripts/Collect-Benchmarks.py"]
    Bench --> BenchGate{"Regression Barrier:<br/>Compare vs BENCHMARK-HISTORY.csv"}
    BenchGate -- "Regression <= 20%" --> DashGen["Generate Visual Dashboard:<br/>benchmarks/BenchmarkCharts.py"]
    BenchGate -- "Regression > 20%" --> Fail["Block Pipeline (CI Failed)"]
    DashGen --> InfraVal["Infrastructure Validation:<br/>integra-infra/scripts/validate.sh (YAML, Helm, K8s, Terraform)"]
```

---

[◀ Back to Chapter](index.md) | [▲ Master Table of Contents](../index.md)


---

[◀ Previous: Global Architecture](../03-global-architecture/index.md) | [▲ Master Table of Contents](../index.md) | [Next: Directory Structure ▶](../05-directory-and-file-structure/index.md)
