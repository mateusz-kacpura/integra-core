# 4. Kluczowe przepływy pracy (System Workflows)

[◀ Poprzedni: Globalna architektura systemu](../03-global-architecture/index.md) | [▲ Główny spis treści](../index.md) | [Następny: Dokumentacja katalogów i plików projektu ▶](../05-directory-and-file-structure/index.md)

---

## 📑 Podrozdziały w tym dziale

- [**4.1. Przepływ 1: Wielopoziomowa weryfikacja tokena JWT i stanu bezpieczeństwa (L1 -> L2 -> L3)**](4.1-jwt-and-security-verification.md)
- [**4.2. Przepływ 2: Atomowe księgowanie podwójnego zapisu (Double-Entry Ledger Posting & Idempotency)**](4.2-double-entry-ledger-posting.md)
- [**4.3. Przepływ 3: Rozproszona saga rejestracji użytkownika i synchronizacji replik odczytu**](4.3-user-registration-saga.md)
- [**4.4. Przepływ 4: Potok jakości CI/CD, audyt bezpieczeństwa i bramka regresji benchmarków**](4.4-cicd-quality-pipeline.md)

---

### 4.1. Przepływ 1: Wielopoziomowa weryfikacja tokena JWT i stanu bezpieczeństwa (L1 -> L2 -> L3)

Poniższy diagram sekwencji ilustruje działanie weryfikatora `auth.JWTVerifier` (`integra-shared-go/pkg/auth/verifier.go`), pamięci podręcznej `L1Cache` (`l1_cache.go`) oraz strażnika stanu konta i ról w `Valkey` (`account_state.go`, `securitystate/authority.go`):

```mermaid
sequenceDiagram
    autonumber
    actor Client as Klient API
    participant MW as Gin Auth Middleware
    participant L1 as L1 In-Process Cache<br/>(32 Shardy, CacheLinePad, TTL 500ms)
    participant Ed as Weryfikator Ed25519<br/>(CPU Crypto)
    participant VK as L2 Valkey Cluster<br/>(ChaCha20-Poly1305 + Lua Lease 1s)
    participant PG as L3 PostgreSQL<br/>(Reconciler Worker)

    Client->>MW: Żądanie HTTP z nagłówkiem Authorization: Bearer <JWT>
    MW->>L1: Pobranie roszczeń (FNV-1a shard hash, 0 alokacji)
    alt L1 CACHE HIT (< 500ms, ~36.5 ns)
        L1-->>MW: Zwrócenie *auth.Claims z pamięci L1 CPU
    else L1 CACHE MISS lub wpis wygasł
        MW->>VK: Sprawdzenie zaszyfrowanych Claims w L2 (XChaCha20-Poly1305)
        alt L2 CLAIMS HIT
            VK-->>MW: Zaszyfrowany blob -> Odszyfrowanie ChaCha20 (~311 ns)
        else L2 CLAIMS MISS
            MW->>Ed: Weryfikacja podpisu asymetrycznego Ed25519 (~60 µs)
            Ed-->>MW: Poprawne *auth.Claims
            MW->>VK: Zapis zaszyfrowanych Claims (ChaCha20-Poly1305)
        end
        MW->>L1: Zapis w lokalnym L1 Cache (TTL = 500ms < 1s Lease)
    end

    MW->>VK: Skrypt Lua accountCheck (HGET {auth:user-status}:snapshot)
    alt Dzierżawa (Lease 1s) aktywna i konto != "blocked"
        VK-->>MW: Status = 1 (OK)
        MW-->>Client: Autoryzacja pomyślna (Kontynuacja żądania)
    else Konto zablokowane (Status = 0)
        VK-->>MW: Status = 0 (ErrUnauthorized)
        MW-->>Client: HTTP 401 Unauthorized
    else Brak dzierżawy / Wygasła (>1s) / Awaria Valkey (Status = -1)
        VK-->>MW: ErrUnavailable (Tryb Fail-Closed!)
        MW-->>Client: HTTP 503 Service Unavailable (Bezpieczne odrzucenie)
        Note over PG,VK: Background Worker cyklicznie odnawia<br/>1-sekundowy Lease z tabel PostgreSQL
    end

```

---

### 4.2. Przepływ 2: Atomowe księgowanie podwójnego zapisu (Double-Entry Ledger Posting & Idempotency)

Przepływ pracy w `integra-ledger-service` podczas wykonywania operacji `PostAccounts` (`internal/infrastructure/postgres/posting.go` oraz `db/queries/posting_commit.sql`):

```mermaid
sequenceDiagram
    autonumber
    actor Caller as Klient / Usługa Zlecająca
    participant Srv as Ledger Service (PostAccounts)
    participant Dom as Domain (Money & PostingHash)
    participant DB as PostgreSQL (Schema: ledger)
    participant Wrk as Outbox & Snapshot Worker
    participant Kafka as Broker Kafka

    Caller->>Srv: RecordTransaction (DebitID, CreditID, Amount, Currency, IdempotencyKey, ExpectedVersions)
    Srv->>Srv: Weryfikacja zgodności TenantID z kontekstem tożsamości
    Srv->>Dom: Obliczenie 33-bajtowego skrótu binarnego RequestDigest v2 (SHA-256, 200.9 ns, 0 B/op)
    Dom-->>Srv: RequestDigest [33]byte

    Srv->>DB: BEGIN TRANSACTION
    Srv->>DB: ClaimAndLockPostingAccounts (Blokada wierszy FOR UPDATE + sprawdzenie IdempotencyKey)
    alt IdempotencyKey już istnieje z tym samym RequestDigest
        DB-->>Srv: Zwrócenie istniejącego TransactionID (Powtórzenie idempotentne)
        Srv->>DB: ROLLBACK
        Srv-->>Caller: 200 OK (Istniejący TransactionID)
    else IdempotencyKey istnieje z INNYM RequestDigest
        DB-->>Srv: Błąd konfliktu parametrów
        Srv->>DB: ROLLBACK
        Srv-->>Caller: 409 Conflict (Naruszenie kontraktu idempotentności)
    else Nowa transakcja
        Srv->>Dom: Walidacja domenowa: Currency match, MinorUnits (int64), Optimistic Locking Versions
        Srv->>DB: Wykonanie pojedynczego zapytania CTE CommitPosting:
        Note right of DB: 1. UPDATE ledger_accounts (Debit & Credit, version + 1)<br/>2. INSERT INTO account_events (historia sald)<br/>3. INSERT INTO snapshot_jobs (jeśli version % interval == 0)<br/>4. INSERT INTO events (Event Store)<br/>5. INSERT INTO ledger_entries (dziennik księgowy)<br/>6. UPDATE posting_requests (Zapis wyniku idempotentnego)<br/>7. INSERT INTO outbox (Zdarzenie z nagłówkami W3C Trace)
        DB-->>Srv: Zatwierdzony wiersz transakcji
        Srv->>DB: COMMIT TRANSACTION
        Srv-->>Caller: 200 OK (TransactionID, RecordedAt)
    end

    par Asynchroniczna publikacja i migawki
        Wrk->>DB: Pobranie nieopublikowanych rekordów z ledger.outbox
        Wrk->>Kafka: Publikacja na temat "ledger-updated" (z propagacją traceparent & correlation_id)
        Wrk->>DB: Oznaczenie published = true w ledger.outbox
    and Budowanie migawek kont
        Wrk->>DB: Pobranie zadań z ledger.snapshot_jobs
        Wrk->>DB: Zapis skonsolidowanego AccountSnapshot (przyspieszenie odczytu 45.2x)
    end

```

---

### 4.3. Przepływ 3: Rozproszona saga rejestracji użytkownika i synchronizacji replik odczytu (Outbox -> Kafka -> Read Replicas)

Przepływ zdarzeń pomiędzy `integra-iam-service`, `integra-user-service` oraz `integra-ledger-service` z zachowaniem pełnego kontekstu śledzenia `OpenTelemetry` (`traceparent`, `correlation_id`):

```mermaid
flowchart LR
    subgraph IAM["1. integra-iam-service"]
        Reg["Rejestracja / Zmiana Roli<br/>(PostgresUnitOfWork)"]
        IAM_DB[("iam.users<br/>iam.user_roles<br/>iam.outbox")]
        IAM_Relay["IAM Outbox Worker"]
    end

    subgraph Broker["2. Apache Kafka"]
        T_User[["Topic: user-events"]]
        T_Role[["Topic: user-roles-sync"]]
        T_Lic[["Topic: license-events"]]
    end

    subgraph Consumers["3. Konsumenci i lokalne repliki odczytu"]
        subgraph UserSvc["integra-user-service"]
            U_Cons["Kafka Consumer<br/>(IdentityHandler)"]
            U_DB[("users.profiles<br/>users.accounts_authz")]
            U_VK["Valkey Barrier State<br/>(Odświeżenie wersji konta)"]
        end

        subgraph LedgerSvc["integra-ledger-service"]
            L_Cons["Kafka Router"]
            L_DB[("ledger.users_read<br/>ledger.orgs_read<br/>(Upsert z kontrolą wersji)")]
        end
    end

    Reg -->|"Transakcja ACID + OutboxHeaders(ctx)"| IAM_DB
    IAM_DB -->|"Polling batch=100"| IAM_Relay
    IAM_Relay --> T_User
    IAM_Relay --> T_Role
    IAM_Relay --> T_Lic

    T_User -->|"telemetry.Restore(headers)"| U_Cons
    T_Role --> U_Cons
    U_Cons --> U_DB
    U_Cons --> U_VK

    T_User -->|"replicas.Upsert (version > current)"| L_Cons
    T_Role --> L_Cons
    T_Lic --> L_Cons
    L_Cons --> L_DB

```

---

### 4.4. Przepływ 4: Potok jakości CI/CD, audyt bezpieczeństwa i bramka regresji benchmarków

Przepływ zdefiniowany w `.github/workflows/ci.yaml`, `scripts/Audit-Project.py`, `scripts/Check-LocalModules.py` oraz `scripts/Collect-Benchmarks.py`:

```mermaid
flowchart TD
    Push["Push / Pull Request do Monorepo"] --> Setup["Inicjalizacja środowiska Go 1.26 (GOWORK=off, CGO_ENABLED=1)<br/>Start kontenerów: 3x Postgres 18, MongoDB 8, Valkey 9"]
    Setup --> Audit["scripts/Audit-Project.py<br/>Weryfikacja każdego modułu w izolacji (go vet, staticcheck v0.7.0)"]
    Audit --> Vuln["govulncheck ./...<br/>Skanowanie podatności CVE w drzewie zależności"]
    Vuln --> Codegen["Weryfikacja generatorów kodu:<br/>check-sqlc.py (spójność SQL) & check-proto.py (spójność Buf Protobuf)"]
    Codegen --> UnitTests["Testy jednostkowe i architektoniczne AST:<br/>go test -race -shuffle=on -mod=readonly ./..."]
    UnitTests --> Migrate["Migracja izolowanych baz testowych:<br/>init-prod-pg.sql + cmd/migrate (IAM, User, Ledger)"]
    Migrate --> IntTests["Testy integracyjne z detektorem wyścigów:<br/>-tags=integration -race (Postgres, Mongo, Valkey, Kafka)"]
    IntTests --> BenchGate{"Bramka Wydajnościowa:<br/>Collect-Benchmarks.py<br/>(--max-regression-pct 20.0)"}
    BenchGate -- "Brak regresji > 20%" --> DashGen["Generowanie wykresów SVG, Mermaid<br/>i BENCHMARKS-DASHBOARD.html"]
    BenchGate -- "Wykryto regresję > 20%" --> Fail["Blokada Pipeline'u (CI Failed)"]
    DashGen --> InfraVal["Walidacja Infrastruktury:<br/>integra-infra/scripts/validate.sh (YAML, Helm, K8s, Terraform)"]

```

---

[◀ Poprzedni: Globalna architektura systemu](../03-global-architecture/index.md) | [▲ Główny spis treści](../index.md) | [Następny: Dokumentacja katalogów i plików projektu ▶](../05-directory-and-file-structure/index.md)
