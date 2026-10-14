# 1. Cel projektu, wartość biznesowa i unikalność na rynku Go

[▲ Główny spis treści](../README.md) | [Następny: Analiza benchmarków i uzasadnienie wydajnościowe ▶](../02-analiza-benchmarkow-i-wydajnosc/README.md)

---

## 📑 Podrozdziały w tym dziale

- [**1.1. Cel projektu i wizja architektoniczna**](1.1-cel-projektu-i-wizja-architektoniczna.md)
- [**1.2. Wartość biznesowa (Business Value)**](1.2-wartosc-biznesowa.md)
- [**1.3. Unikalność w skali rynku projektów w Go**](1.3-unikalnosc-na-rynku-go.md)

---

### 1.1. Cel projektu i wizja architektoniczna

Projekt **Integra Platform** (`[github.com/mateusz-kacpura/integra-enterprise](https://github.com/mateusz-kacpura/integra-enterprise)` oraz wyodrębniony moduł otwartoźródłowy `[github.com/mateusz-kacpura/integra-core](https://github.com/mateusz-kacpura/integra-core)`) to zaawansowany, wielomodułowy ekosystem mikroserwisowy klasy **Enterprise FinTech & B2B SaaS Core**, napisany w **Go 1.26**. System stanowi fundament pod rozbudowaną platformę korporacyjną obejmującą aż **38 kanonicznych strumieni zdarzeń (domen biznesowych)** – od zarządzania tożsamością, uprawnieniami i licencjami (**IAM**), przez profile użytkowników i matryce kompetencji (**User Service**), księgowość podwójnego zapisu, fakturowanie oraz wypłaty (**Ledger Service**), aż po pozostałe moduły systemu.

Projekt realizuje strategię hybrydowej dystrybucji **Open-Core + Enterprise Proprietary**:

* **Warstwa Open-Core (`.dist/integra-core`)**: Wydzielana automatycznie na licencji **Apache 2.0**, dostarczająca publiczne kontrakty gRPC/ConnectRPC (`proto`), sanityzowane logowanie (`zerolog` + `log/slog`), propagację telemetrii W3C (`OpenTelemetry`), szybką serializację JSON (`bytedance/sonic`), paginację kursorową oraz generatory `UUIDv7`.


* **Warstwa Enterprise (`integra-enterprise/*`)**: Objęta licencją komercyjną (*Proprietary and Confidential*), zawierająca krytyczne mechanizmy bezpieczeństwa (weryfikatory Ed25519 z pamięcią podręczną L1 CPU, szyfrowanie `XChaCha20-Poly1305`, strażników licencji i autorytetu ról w `Valkey`), dynamiczny routing SQL (`sqlrouting`) oraz mikroserwisy domenowe (`integra-iam-service`, `integra-user-service`, `integra-ledger-service`) i kompletną infrastrukturę chmurową (`integra-infra`).



---

### 1.2. Wartość biznesowa (Business Value)

1. **Bezkompromisowa integralność finansowa (Double-Entry Bookkeeping & Event Sourcing)**:
* W przeciwieństwie do typowych systemów przechowujących saldo jako pojedynczą kolumnę podatną na błędy wyścigu (*race conditions*), `integra-ledger-service` implementuje pełny **Event Sourcing** połączony z księgowością podwójnego zapisu (*double-entry accounting*), gdzie każde księgowanie wymaga bilansujących się stron **Debit** i **Credit** wyrażonych w niezmiennych jednostkach mniejszych (*minor units / grosze / centy* opartych na dokładnej arytmetyce 64-bitowej `int64` z limitem `MaxMinorUnits = 1 << 50`, eliminującej błędy zaokrągleń IEEE-754).




2. **Natychmiastowe unieważnianie uprawnień i kont bez utraty wydajności (1-Second Fail-Closed Security Lease)**:
* W klasycznych architekturach mikroserwisowych opartych na tokenach JWT istnieje bolesny kompromis: albo tokeny są bezstanowe (i po zablokowaniu konta lub odebraniu roli użytkownik zachowuje dostęp aż do wygaśnięcia tokena, np. przez 15 minut), albo każde żądanie odpytuje bazę danych / Redis, dławiąc przepustowość sieciową.


* Integra rozwiązuje ten problem za pomocą trójpoziomowej architektury **L1 Hardware Cache (TTL 500 ms, max 800 ms) + L2 Valkey Atomic Lua Lease (1000 ms Fail-Closed) + L3 PostgreSQL Reconciler**. Gwarantuje to twardą granicę bezpieczeństwa: zawieszenie konta lub zmiana roli obowiązuje w całym klastrze w czasie **poniżej 1 sekundy**, a w przypadku awarii warstwy cache system przechodzi w bezpieczny tryb **Fail-Closed** (odrzuca dostęp zamiast ufać przeterminowanym danym).




3. **Wbudowana zgodność z RODO / GDPR i PCI-DSS na poziomie strumienia I/O**:
* Pakiet `pkg/logger` implementuje strumieniowy filtr `SanitizingWriter` oraz hooki dla `zerolog` i `log/slog`, które w locie (przed zapisem bajtów na `stdout`, `stderr` lub wysłaniem do kolektora Loki/Fluentd) wykrywają i maskują dane wrażliwe (PII): numery kart kredytowych (formatowane do `4111 **** **** 1234`), numery **PESEL** (`[PESEL_REDACTED]`), numery **NIP** (`[NIP_REDACTED]`), tokeny **JWT** (`[JWT_REDACTED]`) oraz adresy **e-mail** (`u***@domain.com`). Eliminuje to ryzyko wycieku danych osobowych do systemów obserwowalności.




4. **Drastyczna redukcja kosztów chmury obliczeniowej (FinOps)**:
* Dzięki eliminacji alokacji sterty (*Zero-Allocation Hot Paths*) w 60,66% wszystkich testowanych ścieżek krytycznych oraz użyciu pamięci podręcznej L1 i jednozapytaniowych transakcji CTE w PostgreSQL, pojedynczy kontener mikroserwisu w konfiguracji produkcyjnej Helm wymaga zaledwie `100m CPU` i `128Mi RAM` (z limitem `500m CPU` i `512Mi RAM`).





---

### 1.3. Unikalność w skali rynku projektów w Go

W literaturze branżowej Go podkreśla się, że język ten został zaprojektowany przez Google (Roberta Griesemera, Roba Pike'a i Kena Thompsona) jako pragmatyczny język systemowy łączący szybkość wykonywania zbliżoną do C/C++ z produktywnością i automatycznym zarządzaniem pamięcią (garbage collector o pauzach sub-milisekundowych). Jednak większość projektów rynkowych w Go zatrzymuje się na standardowych konwencjach (np. prosty układ `cmd/`, `internal/`, `migrations/`, standardowy `encoding/json`, ORM-y takie jak GORM oraz naiwne użycie muteksów).

Na tle rynku ekosystem **Integra** wyróżnia się rozwiązaniami spotykanymi niemal wyłącznie w silnikach bazodanowych (takich jak CockroachDB czy etcd) oraz systemach High-Frequency Trading:

* **Świadomość mikroarchitektury CPU i eliminacja *false sharing* (`cpu.CacheLinePad`)**: W `integra-user-service/internal/infrastructure/barrier_state.go` oraz `integra-shared-go/pkg/auth/l1_cache.go` struktury współbieżne są wyrównane do 64-bajtowych linii pamięci podręcznej L1D procesora za pomocą `golang.org/x/sys/cpu.CacheLinePad`. Zapobiega to inwalidacji linii cache pomiędzy rdzeniami CPU podczas współbieżnych zapisów atomowych (`atomic.Int64` / `atomic.Uint64`).


* **Projektowanie zorientowane na dane (*Data-Oriented Design* – DOD) i alokacja arenowa**: W `LedgerAggregate` (`integra-ledger-service/internal/domain/aggregate.go`) zamiast rozproszonych alokacji wskaźnikowych na stercie zastosowano prealokowane ciągłe tablice wartości (*arena slices*: `accountArena []AccountState`, `invoiceArena []InvoiceState`, `walletArena []WalletState`, `journalArena []JournalEntryState`). Zwiększa to lokalność przestrzenną danych w pamięci podręcznej procesora i redukuje narzut garbage collectora (GC).


* **Kompilacja zapytań SQL w czasie budowania (`sqlc`) połączona z dynamicznym routingiem L1 (`sqlrouting`)**: Zamiast powolnej refleksji ORM, projekt generuje w 100% typowany kod Go z czystego SQL (`sqlc`), a autorski pakiet `integra-postgres-adapters/sqlrouting` pozwala bezpiecznie podmieniać docelowe tabele (np. `outbox_events` -> `iam.outbox`) z walidacją białej listy (*allowlist*) przeciwko SQL Injection i buforowaniem przepisanych zapytań w pamięci L1 (`19.11 ns/op`, `0 B/op`).


* **Jednozapytaniowy atomowy zapis księgowy (Single-Roundtrip CTE Posting)**: W pliku `posting_commit.sql` cała transakcja księgowa – aktualizacja sald obu kont z kontrolą wersji (*Optimistic Concurrency Control*), zapis historii zmian sald (`account_events`), zlecenie migawki (`snapshot_jobs`), dopisanie zdarzenia do `event_store` (`ledger.events`), utworzenie wpisu w dzienniku (`ledger_entries`), zapis sumy kontrolnej idempotentności (`posting_requests`) oraz publikacja do tabeli `outbox` – wykonuje się w **jednym atomowym zapytaniu SQL CTE (`WITH updated AS (...)`)**, redukując liczbę przebiegów sieciowych (round-trips) do bazy danych do dokładnie jednego.


* **Automatyczne testy architektoniczne AST (`go/parser`)**: W każdym module znajdują się testy `architecture_test.go`, które parsują drzewo składniowe kodu Go (`go/parser`, `go/token`) i wymuszają reguły Clean Architecture (np. zakaz importowania `google.golang.org/grpc` w adapterach współdzielonych, zakaz zależności cyklicznych, zakaz zewnętrznych zależności w rdzeniu `pkg/rbac` poza `sync`, `sync/atomic` i `strings`).

---

[▲ Główny spis treści](../README.md) | [Następny: Analiza benchmarków i uzasadnienie wydajnościowe ▶](../02-analiza-benchmarkow-i-wydajnosc/README.md)
